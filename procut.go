package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func main() {
	algPtr := flag.String("alg", "", "The hashing algorithm (md5, sha1, sha256, sha512).")
	dictPtr := flag.String("dict", "", "The dictionary file.")
	saltPtr := flag.String("salt", "", "The salt to use.")
	outputPtr := flag.String("output", "", "The file to save found passwords to.")
	threadsPtr := flag.Int("threads", runtime.NumCPU(), "Number of threads to use.")
	flag.Parse()

	// Check for the positional argument (hash_input)
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: go run procut.go or ./procut --alg <alg> --dict <dict> [--salt <salt>] [<hash or hash file>] [--output <file>]")
		fmt.Println("\n-> Make sure you use '--' before any argument")
		flag.PrintDefaults()
		os.Exit(1)
	}
	hashInput := flag.Args()[0]

	if *dictPtr == "" || *algPtr == "" {
		fmt.Println("Usage: go run procut.go or ./procut --alg <alg> --dict <dict> [--salt <salt>] [<hash or hash_file>] [--output <file>]")
		fmt.Println("\n-> Make sure you use '--' before any argument")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var hashesToCrack []string
	fileInfo, err := os.Stat(hashInput)
	if err == nil && !fileInfo.IsDir() { // It's a file
		file, err := os.Open(hashInput)
		if err != nil {
			fmt.Printf("-> Error opening hash file '%s': %v\n", hashInput, err)
			os.Exit(1)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			hash := strings.TrimSpace(scanner.Text()) // <--- ADDED TRIMSPACE HERE
			if hash != "" {
				hashesToCrack = append(hashesToCrack, hash)
			}
		}
		if len(hashesToCrack) == 0 {
			fmt.Printf("-> Error: Hash file '%s' is empty or contains no valid hashes.\n", hashInput)
			os.Exit(1)
		}
	} else { // It's a hash string
		hashesToCrack = append(hashesToCrack, hashInput)
	}

	hashFuncFactory := make(map[string]func() hash.Hash)
	hashFuncFactory["md5"] = md5.New
	hashFuncFactory["sha1"] = sha1.New
	hashFuncFactory["sha256"] = sha256.New
	hashFuncFactory["sha512"] = sha512.New

	createHasher, ok := hashFuncFactory[*algPtr]
	if !ok {
		fmt.Printf("-> Unsupported hash algorithm '%s'. \n-> Supported: md5, sha1, sha256, sha512.\n", *algPtr)
		os.Exit(1)
	}

	dictFile, err := os.Open(*dictPtr)
	if err != nil {
		fmt.Printf("-> Error opening dictionary file '%s': %v\n", *dictPtr, err)
		os.Exit(1)
	}
	defer dictFile.Close()

	lines := 0
	scanner := bufio.NewScanner(dictFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) > 0 {
			lines++
		}
	}
	dictFile.Seek(0, 0)

	foundPasswords := make(map[string]string)
	var foundMu sync.Mutex

	for _, targetHash := range hashesToCrack {
		fmt.Printf("\n\nCracking hash: %s\n\n", targetHash)

		passwordsChan := make(chan string, *threadsPtr)
		var wg sync.WaitGroup
		var hashFoundOnce sync.Once
		foundForCurrentHash := false

		bar := progressbar.Default(int64(lines))

		for i := 0; i < *threadsPtr; i++ {
			wg.Add(1)
			go func() {
				h := createHasher()
				defer wg.Done()
				for password := range passwordsChan {
					bar.Add(1)
					if foundForCurrentHash {
						continue
					}
				
				toHash := password + *saltPtr
				h.Reset()
				h.Write([]byte(toHash))
				hashed := hex.EncodeToString(h.Sum(nil))

				if hashed == targetHash {
					hashFoundOnce.Do(func() {
						foundForCurrentHash = true
						fmt.Printf("\n\n> Password found: %s \n-> For hash %s\n", password, targetHash)
						foundMu.Lock()
						foundPasswords[targetHash] = password
						foundMu.Unlock()
						go func() {
							for range passwordsChan {}
						}()
					})
				}
				}
			}()
		}

		dictFile.Seek(0, 0)
		scanner = bufio.NewScanner(dictFile)
		for scanner.Scan() {
			if foundForCurrentHash {
				break
			}
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) > 0 {
				password := parts[len(parts)-1]
				passwordsChan <- password
			}
		}
		close(passwordsChan)
		wg.Wait()

		if !foundForCurrentHash {
			fmt.Println("\n> Password not found in the dictionary.")
		}
	}

	if *outputPtr != "" {
		outFile, err := os.Create(*outputPtr)
		if err != nil {
			fmt.Printf("Error creating output file '%s': %v\n", *outputPtr, err)
			os.Exit(1)
		}
		defer outFile.Close()
		for h, p := range foundPasswords {
			outFile.WriteString(fmt.Sprintf("%s:%s\n", h, p))
		}
		fmt.Printf("\n> Found passwords saved to %s\n", *outputPtr)
	}
}
