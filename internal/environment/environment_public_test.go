package environment_test

import (
	"calabash/internal/environment"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("getting from main environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: nil,
		}

		a := env.Get("a")

		if a != 10 {
			t.Error("did not retrieve value in environment")
		}
	})

	t.Run("getting from parent environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 20,
				},
				Parent: nil,
			},
		}

		b := env.Get("b")

		if b != 20 {
			t.Error("did not retrieve value in parent environment")
		}
	})
}

func TestAdd(t *testing.T) {
	env := &environment.Environment[int]{
		Fields: make(map[string]int),
		Parent: nil,
	}

	if len(env.Fields) != 0 {
		t.Fatal("environment should be empty to start")
	}

	env.Add("int", 30)

	if env.Fields["int"] != 30 {
		t.Error("30 was not set for variable 'int'")
	}

	env.Add("int", 40)

	if env.Fields["int"] != 40 {
		t.Error("40 was not set for variable 'int'")
	}
}

func TestSet(t *testing.T) {
	t.Run("setting in main environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: nil,
		}

		env.Set("a", 20)

		if env.Fields["a"] != 20 {
			t.Error("did not set value in environment")
		}
	})

	t.Run("setting in parent environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 20,
				},
				Parent: nil,
			},
		}

		env.Set("b", 10)

		if env.Parent.Fields["b"] != 10 {
			t.Error("did not set value in parent environment")
		}
	})

	t.Run("nothing is set if variable doesn't exist", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 20,
				},
				Parent: nil,
			},
		}

		env.Set("c", 10)

		_, okA := env.Fields["c"]
		_, okB := env.Parent.Fields["c"]
		if okA || okB {
			t.Error("variable `c` should not have been set")
		}
	})
}

func TestHas(t *testing.T) {
	t.Run("should detect if in main environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: nil,
		}

		ok := env.Has("a")

		if !ok {
			t.Error("did not find variable `a` in environment")
		}
	})

	t.Run("should detect if in parent chain", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		ok := env.Has("c")

		if !ok {
			t.Error("did not find variable `c` in environment chain")
		}
	})

	t.Run("should not detect if not in environment chain", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		ok := env.Has("d")

		if ok {
			t.Error("should not find variable `d` in environment chain")
		}
	})
}

func TestHasDirectly(t *testing.T) {
	t.Run("should not detect if not in main environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		ok := env.HasDirectly("c")

		if ok {
			t.Error("should not find variable `c`")
		}
	})

	t.Run("should detect if in main environment", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		ok := env.HasDirectly("a")

		if !ok {
			t.Error("should find variable `a`")
		}
	})
}

func TestSlice(t *testing.T) {
	t.Run("should receive nil value if slicing 0", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		newEnv := environment.Slice(env, 0)

		if newEnv != nil {
			t.Error("did not receive nil value for slicing 0 layers")
		}
	})

	t.Run("should add one level for value 1", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		newEnv := environment.Slice(env, 1)

		if !reflect.DeepEqual(newEnv, &environment.Environment[int]{Fields: map[string]int{"a": 10}, Parent: nil}) {
			t.Error("First layer of sliced environment should have nil parent")
		}

		if env == newEnv {
			t.Error("Environments should be different pointers")
		}
	})

	t.Run("should add two level for value 2", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		newEnv := environment.Slice(env, 2)

		if !reflect.DeepEqual(newEnv, &environment.Environment[int]{Fields: map[string]int{"a": 10}, Parent: &environment.Environment[int]{Fields: map[string]int{"b": 10}, Parent: nil}}) {
			t.Error("First layer of sliced environment should have nil parent")
		}
	})

	t.Run("should add all levels when value is greater than depth", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		newEnv := environment.Slice(env, 5)

		if !reflect.DeepEqual(newEnv, env) {
			t.Error("First layer of sliced environment should have nil parent")
		}

		if newEnv == env {
			t.Error("Environments should be different pointers")
		}
	})

	t.Run("variables beyond slice depth are unreachable", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: nil,
				},
			},
		}

		newEnv := environment.Slice(env, 2)

		if newEnv.Has("c") {
			t.Error("should not be able to find `c` in sliced environment")
		}
	})

	t.Run("should leave original slice unchanged", func(t *testing.T) {
		env := &environment.Environment[int]{
			Fields: map[string]int{
				"a": 10,
			},
			Parent: &environment.Environment[int]{
				Fields: map[string]int{
					"b": 10,
				},
				Parent: &environment.Environment[int]{
					Fields: map[string]int{
						"c": 10,
					},
					Parent: &environment.Environment[int]{
						Fields: map[string]int{
							"d": 10,
						},
						Parent: nil,
					},
				},
			},
		}

		environment.Slice(env, 2)

		if !env.Has("c") || !env.Has("d") {
			t.Error("should be able to find `c` or `d` in original environment")
		}
	})
}
