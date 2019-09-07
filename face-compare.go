package main

import (
	"flag"
	"fmt"
	"strings"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

var (
	useAll            = flag.Bool("all", false, "Compute all hashes")
	usePHash          = flag.Bool("phash", false, "Compute PHash")
	useAverage        = flag.Bool("average", false, "Compute AverageHash")
	useBlockMean0     = flag.Bool("blockmean0", false, "Compute BlockMeanHash mode 0")
	useBlockMean1     = flag.Bool("blockmean1", false, "Compute BlockMeanHash mode 1")
	useColorMoment    = flag.Bool("colormoment", false, "Compute ColorMomentHash")
	useMarrHildreth   = flag.Bool("marrhildreth", false, "Compute MarrHildrethHash")
	useRadialVariance = flag.Bool("radialvariance", false, "Compute RadialVarianceHash")
)

func setupHashes() []contrib.ImgHashBase {
	var hashes []contrib.ImgHashBase

	if *usePHash || *useAll {
		hashes = append(hashes, contrib.PHash{})
	}
	if *useAverage || *useAll {
		hashes = append(hashes, contrib.AverageHash{})
	}
	if *useBlockMean0 || *useAll {
		hashes = append(hashes, contrib.BlockMeanHash{})
	}
	if *useBlockMean1 || *useAll {
		hashes = append(hashes, contrib.BlockMeanHash{Mode: contrib.BlockMeanHashMode1})
	}
	if *useColorMoment || *useAll {
		hashes = append(hashes, contrib.ColorMomentHash{})
	}
	if *useMarrHildreth || *useAll {
		// MarrHildreth has default parameters for alpha/scale
		hashes = append(hashes, contrib.NewMarrHildrethHash())
	}
	if *useRadialVariance || *useAll {
		// RadialVariance has default parameters too
		hashes = append(hashes, contrib.NewRadialVarianceHash())
	}

	// If no hashes were selected, behave as if all hashes were selected
	if len(hashes) == 0 {
		*useAll = true
		return setupHashes()
	}

	return hashes
}

// FaceCompare is a function to compare 2 face
func FaceCompare(images []gocv.Mat) {
	// var isSamePerson bool
	printHashes := flag.Bool("print", false, "print hash values")
	// construct all of the hash types in a list. normally, you'd only use one of these.
	hashes := setupHashes()
	// compute and compare the images for each hash type
	for _, hash := range hashes {
		results := make([]gocv.Mat, len(images))

		for i, img := range images {
			results[i] = gocv.NewMat()
			defer results[i].Close()
			hash.Compute(img, &results[i])
			if results[i].Empty() {
				fmt.Println("error computing hash for %s")
				return
			}
		}

		// compare for similarity; this returns a float64, but the meaning of values is
		// unique to each algorithm.
		similar := hash.Compare(results[0], results[1])

		// make a pretty name for the hash
		name := strings.TrimPrefix(fmt.Sprintf("%T", hash), "contrib.")
		fmt.Printf("%s: similarity %g\n", name, similar)

		if *printHashes {
			// print hash result for each image
			for i, _ := range images {
				fmt.Printf("%x\n", results[i].ToBytes())
			}
		}
	}
	// return isSamePerson
}
