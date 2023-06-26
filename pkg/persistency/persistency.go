package persistency

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"stats.io/pkg/stats"
)

const persistencyFile string = "./stats.db"

type Pers struct {
	File string
}

func NewPersistent() *Pers {
	p := Pers{
		File: persistencyFile,
	}
	return &p
}

func (p *Pers) LoadFileToMem() *[]stats.TimeEntry {
	file, err := os.Open(p.File)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	// Enumerate over each line
	entries := []stats.TimeEntry{}
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		if len(values) != 2 {
			fmt.Println("Invalid input string")
			return nil
		}

		// Parse the first value as int64
		timeStamp, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing int64 value:", err)
			return nil
		}

		// Parse the second value as uint64
		reqCount, err := strconv.ParseUint(values[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing uint64 value:", err)
			return nil
		}
		entry := stats.TimeEntry{
			TimeStamp: timeStamp,
			Count:     reqCount,
		}
		entries = append(entries, entry)
	}
	return &entries
}

func (p *Pers) DumpToFile(s *stats.Stats) error {

	var sb strings.Builder
	s.Lock.Lock()
	for _, entry := range s.RingBuff {
		entryStr := fmt.Sprintf("%d %d\n", entry.TimeStamp, entry.Count)
		sb.WriteString(entryStr)
	}
	s.Lock.Unlock()
	dataToWrite := sb.String()
	err := ioutil.WriteFile(p.File, []byte(dataToWrite), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}
