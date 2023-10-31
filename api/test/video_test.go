package test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/polldo/govod/core/course"
	"github.com/polldo/govod/core/video"
)

type videoTest struct {
	*TestEnv
}

func TestVideo(t *testing.T) {
	env, err := NewTestEnv(t, "video_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	vt := &videoTest{env}
	ct := &courseTest{env}

	c1 := ct.createCourseOK(t)
	v1 := vt.createVideoOK(t, c1.ID, 1)
	v2 := vt.createVideoOK(t, c1.ID, 2)
	vt.createVideoIndexConflict(t, c1.ID, 2)
	vt.createVideoUnauth(t, c1)

	c2 := ct.createCourseOK(t)
	v3 := vt.createVideoOK(t, c2.ID, 1)
	v3 = vt.updateVideoOK(t, v3)

	vt.showVideoOK(t, v3)
	vs := []video.Video{v1, v2, v3}
	vt.listVideosOK(t, vs)
}

func (vt *videoTest) createVideoOK(t *testing.T, course string, index int) video.Video {
	if err := Login(vt.Server, vt.AdminEmail, vt.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(vt.Server)

	v := video.VideoNew{
		CourseID:    course,
		Index:       index,
		Name:        "Video Test" + strconv.Itoa(rand.Intn(1000)),
		Description: "This is a test video",
		Free:        true,
		URL:         "",
		ImageURL:    "/images/new.png",
	}

	body, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, vt.URL+"/videos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusCreated {
		t.Fatalf("can't create video: status code %s", w.Status)
	}

	var got video.Video
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal created video: %v", err)
	}

	exp := got
	exp.CourseID = v.CourseID
	exp.Index = v.Index
	exp.Name = v.Name
	exp.Description = v.Description
	exp.Free = v.Free
	exp.ImageURL = v.ImageURL

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong video payload. Diff: \n%s", diff)
	}

	return got
}

func (vt *videoTest) createVideoUnauth(t *testing.T, course course.Course) {
	if err := Login(vt.Server, vt.UserEmail, vt.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(vt.Server)

	v := video.VideoNew{
		CourseID:    course.ID,
		Index:       1,
		Name:        "Video Test" + strconv.Itoa(rand.Intn(1000)),
		Description: "This is a test video",
		Free:        true,
		URL:         "",
		ImageURL:    "/images/new.png",
	}

	body, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, vt.URL+"/videos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusUnauthorized {
		t.Fatalf("users should not be able to create videos: status code %s", w.Status)
	}
}

func (vt *videoTest) createVideoIndexConflict(t *testing.T, course string, index int) {
	if err := Login(vt.Server, vt.AdminEmail, vt.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(vt.Server)

	v := video.VideoNew{
		CourseID:    course,
		Index:       index,
		Name:        "Video Test" + strconv.Itoa(rand.Intn(1000)),
		Description: "This is a test video",
		Free:        true,
		URL:         "",
		ImageURL:    "/images/new.png",
	}

	body, err := json.Marshal(&v)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, vt.URL+"/videos", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusBadRequest {
		t.Fatalf("cannot create video with already existing course/index pair: status code %s", w.Status)
	}
}

func (vt *videoTest) updateVideoOK(t *testing.T, v video.Video) video.Video {
	if err := Login(vt.Server, vt.AdminEmail, vt.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(vt.Server)

	vup := video.VideoUp{
		// CourseID:    course,
		Index:       ptr(10),
		Name:        ptr("Video Test" + strconv.Itoa(rand.Intn(1000))),
		Description: ptr("This is a test video"),
		Free:        ptr(true),
		ImageURL:    ptr("/image/updated.png"),
	}

	body, err := json.Marshal(&vup)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPut, vt.URL+"/videos/"+v.ID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't update video: status code %s", w.Status)
	}

	var got video.Video
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal updated video: %v", err)
	}

	exp := got
	exp.Index = *vup.Index
	exp.Name = *vup.Name
	exp.Description = *vup.Description
	exp.Free = *vup.Free
	exp.ImageURL = *vup.ImageURL

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong video payload. Diff: \n%s", diff)
	}
	return got
}

func (vt *videoTest) showVideoOK(t *testing.T, v video.Video) {
	r, err := http.NewRequest(http.MethodGet, vt.URL+"/videos/"+v.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't fetch video: status code %s", w.Status)
	}

	var got video.Video
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal fetched video: %v", err)
	}

	v.CreatedAt = got.CreatedAt
	v.UpdatedAt = got.UpdatedAt

	if diff := cmp.Diff(got, v); diff != "" {
		t.Fatalf("wrong video payload. Diff: \n%s", diff)
	}
}

func (vt *videoTest) listVideosOK(t *testing.T, vs []video.Video) {
	r, err := http.NewRequest(http.MethodGet, vt.URL+"/videos", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := vt.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't fetch videos: status code %s", w.Status)
	}

	var got []video.Video
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal fetched videos: %v", err)
	}

	less := func(a, b video.Video) bool { return a.ID < b.ID }
	if diff := cmp.Diff(got, vs, cmpopts.SortSlices(less)); diff != "" {
		t.Fatalf("wrong videos payload. Diff: \n%s", diff)
	}
}
