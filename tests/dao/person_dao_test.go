package dao_test

import (
	"context"
	"encoding/json"
	"go-postgres-boilerplate/dao"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonDAO(t *testing.T) {
	// Set up the test database
	ctx := context.Background()
	db, cleanup, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer cleanup()

	// Create a new PersonDAO instance
	personDAO, err := dao.NewPersonDAO(db)
	require.NoError(t, err)

	t.Run("CreateAndGetPerson", func(t *testing.T) {
		traits, _ := json.Marshal(map[string]interface{}{"age": 30, "city": "New York"})
		person := &dao.Person{
			UID:       uuid.New().String(),
			Name:      "John Doe",
			Timestamp: time.Now().UTC().Truncate(time.Microsecond),
			Traits:    traits,
		}

		// Create person
		err := personDAO.Save(person)
		assert.NoError(t, err, "failed to create person")

		// Get person
		retrievedPerson, err := personDAO.GetPersonByUID(person.UID)
		assert.NoError(t, err, "failed to get person")
		assert.Equal(t, person, retrievedPerson, "retrieved person does not match expected")
	})

	t.Run("UpdatePerson", func(t *testing.T) {
		initialTraits, _ := json.Marshal(map[string]interface{}{"age": 40, "city": "Chicago"})
		person := &dao.Person{
			UID:       uuid.New().String(),
			Name:      "Bob Smith",
			Timestamp: time.Now().UTC().Truncate(time.Microsecond),
			Traits:    initialTraits,
		}

		// Create person
		err := personDAO.Save(person)
		assert.NoError(t, err, "failed to create person")

		// Update person
		updatedTraits, _ := json.Marshal(map[string]interface{}{"age": 41, "city": "New York"})
		person.Name = "Robert Smith"
		person.Traits = updatedTraits
		person.Timestamp = time.Now().UTC().Truncate(time.Microsecond)

		err = personDAO.Save(person)
		assert.NoError(t, err, "failed to update person")

		// Get updated person
		retrievedPerson, err := personDAO.GetPersonByUID(person.UID)
		assert.NoError(t, err, "failed to get updated person")
		assert.Equal(t, person, retrievedPerson, "updated person does not match expected")
	})

	t.Run("DeletePerson", func(t *testing.T) {
		person := &dao.Person{
			UID:       uuid.New().String(),
			Name:      "Jane Doe",
			Timestamp: time.Now().UTC().Truncate(time.Microsecond),
			Traits:    json.RawMessage(`{"age": 35, "city": "San Francisco"}`),
		}

		// Create person
		err := personDAO.Save(person)
		assert.NoError(t, err, "failed to create person")

		// Delete person
		err = personDAO.DeletePerson(person.UID)
		assert.NoError(t, err, "failed to delete person")

		// Try to get deleted person
		_, err = personDAO.GetPersonByUID(person.UID)
		assert.Equal(t, dao.ErrPersonNotFound, err, "expected ErrPersonNotFound when getting deleted person")
	})

	t.Run("DeletePerson - Not Found", func(t *testing.T) {
		nonExistentUID := uuid.New().String()

		err := personDAO.DeletePerson(nonExistentUID)
		assert.Equal(t, dao.ErrPersonNotFound, err, "expected ErrPersonNotFound when deleting non-existent person")
	})
}
