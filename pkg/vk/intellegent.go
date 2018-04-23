package vk

import (
	"net/url"
	"sort"

	"github.com/zhuharev/vkutil"
)

// GetPhotos get info about user. Returns 5 top (with more likes and
// conain only one face) user photos
func GetPhotos(token string, vkID int) (res []string, err error) {
	u := vkutil.NewWithToken(token)
	u.SetDebug(true)
	vkPhotos, err := u.PhotosGet(vkID, vkutil.PHOTO_PROFILE, url.Values{
		"rev":         {"1"},
		"extended":    {"1"},
		"photo_sizes": {"1"},
		"count":       {"1000"},
	})
	if err != nil {
		return
	}

	photos := Photos(vkPhotos)

	sort.Sort(sort.Reverse(photos))

	for _, photo := range photos {
		url := maxPhotoSize(photo.Sizes)
		// photoImage, err := img.FromURL(url)
		// if err != nil {
		// 	return nil, err
		// }
		// if !img.HasFaces(photoImage) {
		// 	continue
		// }
		res = append(res, url)
	}

	return
}

// Photos implemet sort.Sort interface
type Photos []vkutil.Photo

func (p Photos) Less(i, j int) bool { return p[i].Likes.Count < p[j].Likes.Count }
func (p Photos) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Photos) Len() int           { return len(p) }

func maxPhotoSize(sizes []vkutil.Size) string {
	var (
		max      = ""
		maxWidth = 0
	)
	for _, size := range sizes {
		if size.Width > maxWidth {
			maxWidth = size.Width
			max = size.Source
		}
	}
	return max
}
