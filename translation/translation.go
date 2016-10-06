// this package provides the function for retriving villages locations, borders as well as 
// translation 
package translation

import (
	convexhull "github.com/thomaspeugeot/go-convexhull/convexhull"
)


type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window (used to compute liste of villages)
	sourceCountry Country
	targetCountry Country
}

func (t * Translation) Init(sourceCountry, targetCountry Country) {

	Info.Printf("Init : Source Country is %s with step %d", sourceCountry.Name, sourceCountry.Step)
	Info.Printf("Init : Target Country is %s with step %d", sourceCountry.Name, sourceCountry.Step)

	t.sourceCountry = sourceCountry
	t.sourceCountry.Init()

	t.targetCountry = targetCountry
	t.targetCountry.Init()

}

func (t * Translation) SetRenderingWindow( xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}


// 
func (t * Translation) VillageCoordinates( lat, lng float64) (x, y, distance, latClosest, lngClosest, xSpread, ySpread float64, closestIndex int) {

	// convert from lat lng to x, y in the Country 
	return t.sourceCountry.VillageCoordinates( lat, lng)
}

// from a coordinate in source coutry, get closest body, compute
func (t * Translation) TargetVillage( xSpread, ySpread float64) (latTarget, lngTarget float64) {

	Info.Printf("TargetVillage input xSpread %f ySpread %f", xSpread, ySpread)

	latTarget, lngTarget = t.targetCountry.XYSpreadToLatLngOrig( xSpread, ySpread)

	Info.Printf("TargetVillage output lat %f lng %f", latTarget, lngTarget)

	return latTarget, lngTarget
}

// from a coordinate in source coutry, get border
func (t * Translation) TargetBorder( xSpread, ySpread float64) convexhull.PointList {

	Info.Printf("TargetBorder input xSpread %f ySpread %f", xSpread, ySpread)

	points := t.targetCountry.XYSpreadToLatLngOrigVillage( xSpread, ySpread)

	return points
}

func (t * Translation) SourceBorder( lat, lng float64) convexhull.PointList {

	points := t.sourceCountry.VillageBorder( (t.sourceCountry).LatLng2XY( lat, lng))

	return points
}






