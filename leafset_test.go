package pastry

import (
	"testing"
)

// Test insertion of a node into the routing table
func TestLeafSetInsert(t *testing.T) {
	self_id, err := NodeIDFromBytes([]byte("this is a test Node for testing purposes only."))
	if err != nil {
		t.Fatalf(err.Error())
	}
	self := NewNode(self_id, "127.0.0.1", "127.0.0.1", "testing", 55555)
	t.Logf("%s\n", self_id.String())

	other_id, err := NodeIDFromBytes([]byte("this is some other Node for testing purposes only."))
	if err != nil {
		t.Fatalf(err.Error())
	}
	other := NewNode(other_id, "127.0.0.2", "127.0.0.2", "testing", 55555)
	t.Logf("%s\n", other_id.String())
	t.Logf("Diff: %v\n", self_id.Diff(other_id))
	leafset := NewLeafSet(self)
	go leafset.listen()
	defer leafset.Stop()
	r, err := leafset.Insert(other)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r == nil {
		t.Fatalf("Nil response returned.")
	}
	if r.Pos != 0 {
		t.Fatalf("Expected node to be in pos %d, was put in %d instead.", 0, r.Pos)
	}
	side := self_id.RelPos(other_id)
	side_str := "left"
	other_side_str := "right"
	if side == 1 {
		side_str = "right"
		other_side_str = "left"
	}
	if (r.Left && side == 1) || (!r.Left && side == -1) {
		t.Fatalf("Expected node to be to the %s, was to the %s instead.", side_str, other_side_str)
	}
	r2, err := leafset.Get(nil, 0, r.Left)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r2 == nil {
		t.Fatalf("Nil response returned.")
	}
	if r2.Node == nil {
		t.Fatalf("Expected node, got nil instead.")
	}
	if !r2.Node.ID.Equals(other_id) {
		t.Fatalf("Expected Node %s, got Node %s instead.", other_id, r2.Node.ID)
	}
}

// Test handling of a Node being inserted twice.
func TestLeafSetDoubleInsert(t *testing.T) {
	self_id, err := NodeIDFromBytes([]byte("this is a test Node for testing purposes only."))
	if err != nil {
		t.Fatalf(err.Error())
	}
	self := NewNode(self_id, "127.0.0.1", "127.0.0.1", "testing", 55555)

	other_id, err := NodeIDFromBytes([]byte("this is some other Node for testing purposes only."))
	if err != nil {
		t.Fatalf(err.Error())
	}
	other := NewNode(other_id, "127.0.0.2", "127.0.0.2", "testing", 55555)
	t.Log(other)
	leafset := NewLeafSet(self)
	go leafset.listen()
	defer leafset.Stop()
	r, err := leafset.Insert(other)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r == nil {
		t.Fatalf("First insert returned a nil response.")
	}
	r2, err := leafset.Insert(other)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if r2 == nil {
		t.Fatalf("Second insert returned a nil response.")
	}
	if r.Pos != r2.Pos {
		t.Errorf("Positions expected to be equal. %d != %d", r.Pos, r2.Pos)
	}
}

// Test retrieving a Node by ID
func TestLeafSetGetByID(t *testing.T) {
	self_id, err := NodeIDFromBytes([]byte("this is a test Node for testing purposes only."))
	if err != nil {
		t.Fatal(err.Error())
	}
	self := NewNode(self_id, "127.0.0.1", "127.0.0.1", "testing", 55555)

	leafset := NewLeafSet(self)
	go leafset.listen()
	defer leafset.Stop()

	other_id, err := NodeIDFromBytes([]byte("this is some other Node for testing purposes only."))
	if err != nil {
		t.Fatal(err.Error())
	}
	other := NewNode(other_id, "127.0.0.2", "127.0.0.2", "testing", 55555)
	r, err := leafset.Insert(other)
	if err != nil {
		t.Fatal(err.Error())
	}
	if r == nil {
		t.Fatal("Insert returned nil response.")
	}
	r2, err := leafset.Get(other, -1, false)
	if err != nil {
		t.Fatal(err.Error())
	}
	if r2 == nil {
		t.Fatal("Returned nil response.")
	}
	if r2.Pos != r.Pos {
		t.Errorf("Expected pos %v, got pos %v.", r.Pos, r2.Pos)
	}
	if r2.Left != r.Left {
		expectation := "left"
		result := "right"
		if !r.Left {
			expectation = "right"
			result = "left"
		}
		t.Errorf("Expected node to be on the %v, but it was inserted on the %v.", expectation, result)
	}
	if r2.Node == nil {
		t.Fatalf("r2 returned nil node")
	}
	if r.Node == nil {
		t.Fatalf("r returned nil node")
	}
	if !r2.Node.ID.Equals(r.Node.ID) {
		t.Errorf("Expected node %v, got node %v.", r.Node.ID, r2.Node.ID)
	}
}

// Test retrieving a node by position
func TestLeafSetGetByPos(t *testing.T) {
	self_id, err := NodeIDFromBytes([]byte("this is a test Node for testing purposes only."))
	if err != nil {
		t.Fatal(err.Error())
	}
	self := NewNode(self_id, "127.0.0.1", "127.0.0.1", "testing", 55555)

	leafset := NewLeafSet(self)
	go leafset.listen()
	defer leafset.Stop()

	other_id, err := NodeIDFromBytes([]byte("This is another test Node for testing purposes only."))
	if err != nil {
		t.Fatal(err.Error())
	}
	other := NewNode(other_id, "127.0.0.2", "127.0.0.2", "testing", 55555)
	r, err := leafset.Insert(other)
	if err != nil {
		t.Fatal(err.Error())
	}
	r2, err := leafset.Get(nil, r.Pos, r.Left)
	if err != nil {
		t.Fatal(err.Error())
	}
	if r2 == nil {
		t.Fatal("Returned nil response.")
	}
	if r2.Node == nil {
		t.Fatalf("r2 returned nil node")
	}
	if r.Node == nil {
		t.Fatalf("r returned nil node")
	}
	if !r2.Node.ID.Equals(r.Node.ID) {
		t.Errorf("Expected node %v, got node %v.", r.Node.ID, r2.Node.ID)
	}
}
