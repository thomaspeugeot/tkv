//
// grump reader parses a country specific file and generate a config file of bodies
//
// for each cell of the country specifc file, this program generate bodies per cells according to
// the population count in the cell
//
// the arrangement of circle in each cell is taken from a outsite source (csq something) up to 200 circles 
//
// usage grump-reader -country=xxx where xxx isthe 3 small letter ISO 3166 code for the country (for instance "fra")
// 
package main

import "flag"
import "math"
// import "math/rand"
import "fmt"
import "os"
import "log"
import "bufio"
import "path/filepath"
import "github.com/thomaspeugeot/tkv/barnes-hut"
import "github.com/thomaspeugeot/tkv/quadtree"

// store country code
type country struct {
	Name string
	NCols, NRows, XllCorner, YllCorner int
}	

// coordinates of arrangement of circle packing in a square
type circleCoord struct {
	x,y float64
}

var targetMaxBodies = 400000

var maxCirclePerCell = 500

// storage of circle arrangement per number of circle in the square
type arrangementsStore [][]circleCoord



//
// on the PC
//  go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data"
func main() {

	// flag "country"
	countryPtr := flag.String("country","fra","iso 3166 country code")

	// get the directory containing tkv data through the flag "tkvdata"
	dirTKVDataPtr := flag.String("tkvdata","/Users/thomaspeugeot/the-mapping-data/%s_grumpv1_pcount_00_ascii_30/","directory containing input tkv data")
		
	var country country

	flag.Parse()
	fmt.Println( "country to parse", *countryPtr)
	country.Name = *countryPtr
	fmt.Println( "directory containing tkv data", *dirTKVDataPtr)
	dirTKVData := *dirTKVDataPtr

	// create the path to the agragate country count
	grumpFilePath := fmt.Sprintf( "%s/%s_grumpv1_pcount_00_ascii_30/%sup00ag.asc", dirTKVData, *countryPtr, *countryPtr )
	fmt.Println("relative path ", filepath.Clean( grumpFilePath))
	var grumpFile *os.File
	var err error
	grumpFile, err = os.Open( filepath.Clean( grumpFilePath))
	if err != nil {
		log.Fatal(err)
	}	

	// parse the grump
	var word int
	scanner := bufio.NewScanner( grumpFile)
	scanner.Split(bufio.ScanWords)

	// scan 8 first lines
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.NCols)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.NRows)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.XllCorner)
	scanner.Scan(); scanner.Scan()
	fmt.Sscanf( scanner.Text(), "%d", & country.YllCorner)

	fmt.Println( country )

	// scan the reamining header
	for word < 4 {
		scanner.Scan()
		word++		
		fmt.Println( fmt.Sprintf("item %d : %s", word, scanner.Text()))
	}
	rowLatWidth := 0.0083333333333
	colLngWidth := 0.0083333333333

	// prepare the count matrix
	countMatrix := make([]float64, country.NRows * country.NCols)

	popTotal := 0.0
	// scan the file and store result in countMatrix
	for row :=0; row < country.NRows; row++ {
		lat := float64( country.YllCorner) + (float64( country.NRows - row)*rowLatWidth)
		for col :=0; col < country.NCols ; col++ {
			scanner.Scan()
			// lng := float64(country.XllCorner) + (float64(col)*colLngWidth)

			var count float64
			fmt.Sscanf( scanner.Text(), "%f", &count)
			popTotal += count

			countMatrix[ (country.NRows-row-1)*country.NCols + col ] = count
		}
		fmt.Printf("\rrow %5d lat %2.3f total %f", row, lat, popTotal)
	}
	fmt.Println("")
	grumpFile.Close()

	// get the arrangement
	arrangements := make( arrangementsStore, maxCirclePerCell)
	for nbCircles := 1; nbCircles < maxCirclePerCell; nbCircles++ {

		fmt.Printf("\rgetting arrangement for %3d circles", nbCircles)

		arrangements[nbCircles] = make( []circleCoord, nbCircles)
		

		
		// open the reference file
		circlePackingFilePath := fmt.Sprintf( "%s/csq_coords/csq%d.txt", dirTKVData, nbCircles )
		var circlePackingFile *os.File
		var errCirclePackingFile error
		circlePackingFile, errCirclePackingFile = os.Open( filepath.Clean( circlePackingFilePath))
		if errCirclePackingFile != nil {
			log.Fatal(err)
		}	

		// prepare scanner
		scannerCircle := bufio.NewScanner( circlePackingFile)
		scannerCircle.Split(bufio.ScanWords)

		// one line per circle
		for circle := 0; circle < nbCircles; circle++ {
			
			// scan the id of the circle
			scannerCircle.Scan(); 

			// scan X coordinate
			scannerCircle.Scan()
			fmt.Sscanf( scannerCircle.Text(), "%f", & (arrangements[nbCircles][circle].x))
			// scan Y coordinate
			scannerCircle.Scan()
			fmt.Sscanf( scannerCircle.Text(), "%f", & (arrangements[nbCircles][circle].y))
			// fmt.Printf("getting arrangement for %d circle %f %f\n", nbCircles, arrangements[nbCircles][circle].x, arrangements[nbCircles][circle].y)
		}
		circlePackingFile.Close()
	}

	// prepare the output density file
	var bodies []quadtree.Body
	bodiesInCellMax := 0

	cumulativePopTotal := 0.0
	bodiesNb :=0
	for row :=0; row < country.NRows; row++ {
		lat := float64( country.YllCorner) + (float64( country.NRows - row)*rowLatWidth)
		for col :=0; col < country.NCols ; col++ {
			lng := float64(country.XllCorner) + (float64(col)*colLngWidth)

			// compute relative coordinate of the cell
			relX := (lng - float64(country.XllCorner)) / (float64(country.NCols) * colLngWidth)
			relY := (lat - float64(country.YllCorner)) / (float64(country.NRows) * rowLatWidth)

			// fetch count of the cell
			count := countMatrix[ row*country.NCols + col ]

			// how many bodies ? it is maxBodies *( count / country.PCount) 
			bodiesInCell := int( math.Floor( float64( targetMaxBodies) * (count/popTotal)))
			if bodiesInCell > bodiesInCellMax { bodiesInCellMax = bodiesInCell}
			
			// initiate the bodies
			for i :=0; i<bodiesInCell; i++ {
				var body quadtree.Body
				// angle := float64(i) * 2.0 * math.Pi / float64(bodiesInCell)
				body.X = relX + (1.0/float64(country.NCols))*(0.5 + arrangements[bodiesInCell][i].x)
				body.Y = relY + (1.0/float64(country.NRows))*(0.5 + arrangements[bodiesInCell][i].y)
				body.M = count/float64(bodiesInCell)
				bodies = append( bodies,  body)
			}
			cumulativePopTotal += count
			bodiesNb += bodiesInCell
		}
	}

	// var quadtree quadtree.Quadtree
	// quadtree.Init( &bodies)
	fmt.Println("bodies in cell max ", bodiesInCellMax)
	fmt.Println("cumulative pop ", cumulativePopTotal)
	fmt.Println("nb of bodies ", bodiesNb)

	var run barnes_hut.Run
	run.Init( & bodies)

	run.CaptureConfigCountry( country.Name)
}