package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/*
	Prints edges to a file for visualization with nx and matplotlib in python
*/

func treeOutput(tree []Edge, outfile string) {
	fo, err := os.Create(outfile)
	if err != nil {
		panic(err)
	}

	defer fo.Close()

	writer := bufio.NewWriter(fo)

	var lines []string
	for i := range tree {
		coords1x := tree[i].origin.x
		coords1y := tree[i].origin.y
		coords2x := tree[i].dest.x
		coords2y := tree[i].dest.y
		coordLine := strconv.FormatFloat(coords1x, 'f', 2, 64) + " " + strconv.FormatFloat(coords1y, 'f', 2, 64) + " " + strconv.FormatFloat(coords2x, 'f', 2, 64) + " " + strconv.FormatFloat(coords2y, 'f', 2, 64)
		lines = append(lines, coordLine)
	}

	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Flush the writer to ensure all data is written
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Lines written to file successfully.")
}
