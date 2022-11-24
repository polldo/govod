package test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"testing"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/polldo/govod/core/course"
)

type courseTest struct {
	*TestEnv
}

func TestCourse(t *testing.T) {
	env, err := NewTestEnv(t, "course_test")
	if err != nil {
		t.Fatalf("initializing test env: %v", err)
	}

	ct := &courseTest{env}
	c1 := ct.createCourseOK(t)
	c2 := ct.createCourseOK(t)
	c2 = ct.updateCourseOK(t, c2)
	cs := []course.Course{c1, c2}

	ct.createCourseUnauth(t)
	ct.updateCourseUnauth(t, c2)

	ct.showCourseOK(t, c1)
	ct.showCourseInvalid(t)
	ct.showCourseNotFound(t)
	ct.listCoursesOK(t, cs)
}

func (ct *courseTest) createCourseOK(t *testing.T) course.Course {
	if err := Login(ct.Server, ct.AdminEmail, ct.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := course.CourseNew{
		Name:        "Test" + strconv.Itoa(rand.Intn(1000)),
		Description: "This is a test course",
		Price:       rand.Float64() * 100,
	}

	body, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, ct.URL+"/courses", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusCreated {
		t.Fatalf("can't create course: status code %s", w.Status)
	}

	var got course.Course
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal created course: %v", err)
	}

	exp := got
	exp.Name = c.Name
	exp.Description = c.Description
	exp.Price = c.Price

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong course payload. Diff: \n%s", diff)
	}

	return got
}

func (ct *courseTest) createCourseUnauth(t *testing.T) {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := course.CourseNew{
		Name:        "Test",
		Description: "This is a test course",
		Price:       100,
	}

	body, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, ct.URL+"/courses", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusUnauthorized {
		t.Fatalf("users should not be able to create courses: status code %s", w.Status)
	}
}

func ptr[T any](a T) *T {
	return &a
}

func (ct *courseTest) updateCourseOK(t *testing.T, crs course.Course) course.Course {
	if err := Login(ct.Server, ct.AdminEmail, ct.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := course.CourseUp{
		Name:        ptr("Updated Test"),
		Description: ptr("This is an updated test course"),
		Price:       ptr(500.0),
	}

	body, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPut, ct.URL+"/courses/"+crs.ID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't update course: status code %s", w.Status)
	}

	var got course.Course
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal updated course: %v", err)
	}

	exp := got
	exp.Name = *c.Name
	exp.Description = *c.Description
	exp.Price = *c.Price

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("wrong course payload. Diff: \n%s", diff)
	}
	return got
}

func (ct *courseTest) updateCourseUnauth(t *testing.T, crs course.Course) {
	if err := Login(ct.Server, ct.UserEmail, ct.UserPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := course.CourseUp{
		Name:        ptr("Updated Test Unauth"),
		Description: ptr("This is an updated test course - unauth"),
		Price:       ptr(300.0),
	}

	body, err := json.Marshal(&c)
	if err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPut, ct.URL+"/courses/"+crs.ID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusUnauthorized {
		t.Fatalf("users should not be able to update courses: status code %s", w.Status)
	}
}

func (ct *courseTest) showCourseOK(t *testing.T, crs course.Course) {
	r, err := http.NewRequest(http.MethodGet, ct.URL+"/courses/"+crs.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't fetch course: status code %s", w.Status)
	}

	var got course.Course
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal fetched course: %v", err)
	}

	if diff := cmp.Diff(got, crs); diff != "" {
		t.Fatalf("wrong course payload. Diff: \n%s", diff)
	}
}

func (ct *courseTest) showCourseInvalid(t *testing.T) {
	invalidID := "SELECT * FROM Users WHERE ((Username='$username') AND (Password=MD5('$password')))"

	r, err := http.NewRequest(http.MethodGet, ct.URL+"/courses/"+invalidID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode == http.StatusOK {
		t.Fatalf("ill requests should be blocked: status code %s", w.Status)
	}
}

func (ct *courseTest) showCourseNotFound(t *testing.T) {
	randomID := validate.GenerateID()

	r, err := http.NewRequest(http.MethodGet, ct.URL+"/courses/"+randomID, nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode == http.StatusNotFound {
		t.Fatalf("fetche course should not exist: status code %s", w.Status)
	}
}

func (ct *courseTest) listCoursesOK(t *testing.T, crs []course.Course) {
	r, err := http.NewRequest(http.MethodGet, ct.URL+"/courses", nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := ct.Client().Do(r)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Fatalf("can't fetch course: status code %s", w.Status)
	}

	var got []course.Course
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("cannot unmarshal fetched courses: %v", err)
	}

	less := func(a, b course.Course) bool { return a.ID < b.ID }
	if diff := cmp.Diff(got, crs, cmpopts.SortSlices(less)); diff != "" {
		t.Fatalf("wrong courses payload. Diff: \n%s", diff)
	}
}
