package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"tfinder/color"
	"tfinder/config"
)

type TodoItem struct {
	Path    string
	Line    int
	Pattern string
	Comment string
}

// BUG: this is a bug
var commentPatterns = []*regexp.Regexp{
	// Single line comments
	regexp.MustCompile(`^\s*\/\/(.+)$`), // C, C++, Go, Java, JavaScript, etc.
	regexp.MustCompile(`^\s*#(.+)$`),    // Python, Ruby, Bash, etc.
	regexp.MustCompile(`^\s*--(.+)$`),   // SQL, Lua, etc.
	regexp.MustCompile(`^\s*%(.+)$`),    // Matlab, LaTeX, etc.
	regexp.MustCompile(`^\s*;(.+)$`),    // Assembly, Lisp, etc.

	// Multi-line comment starters
	regexp.MustCompile(`^\s*\/\*+(.+)$`),  // C, C++, Java, etc. start
	regexp.MustCompile(`^\s*\*(.+)$`),     // C, C++, Java, etc. middle
	regexp.MustCompile(`^\s*\*+\/(.+)?$`), // C, C++, Java, etc. end

	// HTML/XML comments
	regexp.MustCompile(`^\s*<!--(.+)$`), // HTML start
	regexp.MustCompile(`^\s*-->(.+)?$`), // HTML end
}

func main() {
	//config := LoadConfig()
	cfg := config.LoadConfig()

	dir := flag.String("dir", cfg.Dir, "The path to look for comments")
	ignorePatterns := flag.String("ignore", strings.Join(cfg.Ignore, ","), "Comma Seperated list of patterns")
	pattern := flag.String("pattern", strings.Join(cfg.Pattern, ","), "Pattern To Search")
	createConfig := flag.Bool("create-config", false, "Create default config")
	onlyComments := flag.Bool("only-comments", true, "Only search in comments")

	flag.Parse()

	if *createConfig {
		err := config.CreateConfig()
		if err != nil {
			fmt.Printf("Error creating config file\n")
			os.Exit(1)
		}
		fmt.Printf("Created default config file\n")
		return
	}

	ignoreList := strings.Split(*ignorePatterns, ",")

	patternList := strings.Split(*pattern, ",")

	patternRegexps := make([]*regexp.Regexp, 0, len(patternList))
	for _, pattern := range patternList {
		if pattern == "" {
			continue
		}
		r, err := regexp.Compile(`(?i)\b` + regexp.QuoteMeta(pattern) + `\b`)
		if err != nil {
			fmt.Printf("Warning: Invalid pattern '%s': %v\n", pattern, err)
			continue
		}
		patternRegexps = append(patternRegexps, r)
	}

	if len(patternRegexps) == 0 {
		fmt.Println("no valid search patterns provided")
		os.Exit(1)
	}

	items := []TodoItem{}

	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path (%s): %v\n", path, err)
			return nil
		}

		if info.IsDir() {
			for _, pattern := range ignoreList {
				if matched, _ := filepath.Match(pattern, info.Name()); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}

		for _, pattern := range ignoreList {
			if strings.Contains(path, pattern) {
				return nil
			}
		}

		ext := strings.ToLower(filepath.Ext(path))
		if likelyBinaryFile(ext) {
			fmt.Printf("Binary file found %s\n", path)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", path, err)
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			if *onlyComments {
				var commentText string
				isComment := false

				for _, commentPattern := range commentPatterns {
					if matches := commentPattern.FindStringSubmatch(line); matches != nil {
						if len(matches) > 1 {
							commentText = matches[1]
							isComment = true
							break
						}
					}
				}

				if !isComment {
					continue
				}

				for i, regex := range patternRegexps {
					if regex.MatchString(commentText) {
						items = append(items, TodoItem{
							Path:    path,
							Line:    lineNum,
							Pattern: patternList[i],
							Comment: strings.TrimSpace(commentText),
						})
						break
					}
				}
			} else {
				for i, regex := range patternRegexps {
					if regex.MatchString(line) {
						items = append(items, TodoItem{
							Path:    path,
							Line:    lineNum,
							Pattern: patternList[i],
							Comment: strings.TrimSpace(line),
						})
						break
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", path, err)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking dirs %v\n", err)
		os.Exit(1)
	}

	if len(items) == 0 {
		fmt.Println("No Item Found")
		return
	}

	fmt.Println("Found some items ->")
	for i, item := range items {
		fmt.Printf("%d. [%s] %s (line %d): \n%s\n\n",
			i+1,
			color.Colorize(item.Pattern, item.Pattern, cfg),
			item.Path,
			item.Line,
			item.Comment,
		)
	}
}

func likelyBinaryFile(ext string) bool {
	binaryExts := map[string]bool{
		".exe": true, ".bin": true, ".obj": true, ".o": true,
		".dll": true, ".png": true, ".jpg": true, ".jpeg": true,
		".gif": true, ".zip": true, ".tar": true, ".gz": true,
		".pdf": true, ".doc": true, ".docx": true, ".pdg": true,
	}
	return binaryExts[ext]
}
