package object

import "testing"

func TestStringHashKey(t *testing.T) {
    hello1 := &String{Value: "Hello World"}
    hello2 := &String{Value: "Hello World"}
    diff1 := &String{Value: "My name is johnny"}
    diff2 := &String{Value: "My name is johnny"}

    if hello1.GenHashKey() != hello2.GenHashKey() {
        t.Errorf("strings with same content have different hash keys")
    }

    if diff1.GenHashKey() != diff2.GenHashKey() {
        t.Errorf("strings with same content have different hash keys")
    }

    if hello1.GenHashKey() == diff1.GenHashKey() {
        t.Errorf("strings with different content have same hash keys")
    }
}

func TestBooleanHashKey(t *testing.T) {
    true1 := &Boolean{Value: true}
    true2 := &Boolean{Value: true}
    false1 := &Boolean{Value: false}
    false2 := &Boolean{Value: false}

    if true1.GenHashKey() != true2.GenHashKey() {
        t.Errorf("trues do not have same hash key")
    }

    if false1.GenHashKey() != false2.GenHashKey() {
        t.Errorf("falses do not have same hash key")
    }

    if true1.GenHashKey() == false1.GenHashKey() {
        t.Errorf("true has same hash key as false")
    }
}

func TestIntegerHashKey(t *testing.T) {
    one1 := &Integer{Value: 1}
    one2 := &Integer{Value: 1}
    two1 := &Integer{Value: 2}
    two2 := &Integer{Value: 2}

    if one1.GenHashKey() != one2.GenHashKey() {
        t.Errorf("integers with same content have twoerent hash keys")
    }

    if two1.GenHashKey() != two2.GenHashKey() {
        t.Errorf("integers with same content have twoerent hash keys")
    }

    if one1.GenHashKey() == two1.GenHashKey() {
        t.Errorf("integers with twoerent content have same hash keys")
    }
}