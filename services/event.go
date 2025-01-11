package services

import (
	"github.com/yasharya2901/smart_divide/models"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (ec *EventService) CreateEvent(name string) (*models.Event, error) {
	// Create an event
	event := models.Event{
		Name: name,
	}
	if err := ec.db.Create(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (ec *EventService) GetEvents() ([]models.Event, error) {
	// Get all events
	var events []models.Event
	if err := ec.db.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (ec *EventService) GetEventByID(id uint, preloadPerson bool) (*models.Event, error) {
	// Get an event by ID
	var event models.Event
	if preloadPerson {
		if err := ec.db.Preload("People").First(&event, id).Error; err != nil {
			return nil, err
		}
	} else {
		if err := ec.db.First(&event, id).Error; err != nil {
			return nil, err
		}
	}
	return &event, nil
}

func (ec *EventService) UpdateEvent(id uint, name string) (*models.Event, error) {
	// Update an event
	var event models.Event
	if err := ec.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	event.Name = name
	if err := ec.db.Save(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (ec *EventService) DeleteEvent(id uint) error {
	// Delete an event
	var event models.Event
	if err := ec.db.First(&event, id).Error; err != nil {
		return err
	}
	if err := ec.db.Delete(&event).Error; err != nil {
		return err
	}

	return nil
}

func (ec *EventService) AddPersonToEvent(eventID, personID uint) error {
	// Add a person to an event
	var event models.Event

	if err := ec.db.Preload("People").First(&event, eventID).Error; err != nil {
		return err
	}

	var person models.Person
	if err := ec.db.First(&person, personID).Error; err != nil {
		return err
	}

	event.People = append(event.People, person)

	if err := ec.db.Save(&event).Error; err != nil {
		return err
	}

	return nil
}

func (ec *EventService) RemovePersonFromEvent(eventID, personID uint) error {
	// Remove a person from an event
	var event models.Event

	if err := ec.db.Preload("People").First(&event, eventID).Error; err != nil {
		return err
	}

	var person models.Person
	if err := ec.db.First(&person, personID).Error; err != nil {
		return err
	}

	for i, p := range event.People {
		if p.ID == person.ID {
			event.People = append(event.People[:i], event.People[i+1:]...)
			break
		}
	}

	if err := ec.db.Save(&event).Error; err != nil {
		return err
	}

	return nil
}
