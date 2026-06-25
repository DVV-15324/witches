// package generator
package generator

// GenMarker là interface đánh dấu struct cần sinh code
type GenMarker interface {
	GenEasyJson() // method rỗng, chỉ để đánh dấu
}
