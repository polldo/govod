package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	ut := &courseTest{env}
	ut.createCourseOK(t)
	ut.createCourseUnauth(t)
}

func (ct *courseTest) createCourseOK(t *testing.T) course.Course {
	if err := Login(ct.Server, ct.AdminEmail, ct.AdminPass); err != nil {
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

func (ct *courseTest) updateCourseOK(t *testing.T, crs course.Course) {
	if err := Login(ct.Server, ct.AdminEmail, ct.AdminPass); err != nil {
		t.Fatal(err)
	}
	defer Logout(ct.Server)

	c := course.CourseUp{
		Name:        ptr("Updated Test"),
		Description: ptr("This is an updated test course"),
		Price:       ptr(200.0),
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

func (ct *courseTest) listCoursesOK(t *testing.T, id string) {
}

func (ct *courseTest) showCourseOK(t *testing.T, id string) {
}

func (ct *courseTest) showCourseInvalid(t *testing.T, id string) {
}
