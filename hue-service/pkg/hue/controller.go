package hue

import (
	"fmt"
	"github.com/amimof/huego"
	"github.com/go-pg/pg/v10"
)

type Controller struct {
	bridges map[string]*huego.Bridge
	db      *pg.DB
}

func loadBridges(db *pg.DB) (map[string]*huego.Bridge, error) {
	var bridges []Bridge
	if err := db.Model(&bridges).Select(); err != nil {
		return nil, err
	}

	huegoBridges := make(map[string]*huego.Bridge)
	for _, bridge := range bridges {
		huegoBridge := huego.New(bridge.Url, bridge.Username)
		huegoBridges[huegoBridge.ID] = huegoBridge
	}
	return huegoBridges, nil
}

func New(db *pg.DB) *Controller {
	bridges, err := loadBridges(db)
	if err != nil {
		panic(err)
	}
	return &Controller{
		bridges: bridges,
		db:      db,
	}
}

func (c *Controller) Bridge(id string) (*huego.Bridge, error) {
	bridge, exists := c.bridges[id]
	if !exists {
		return nil, fmt.Errorf("bridge with id: %s does not exist", id)
	}
	return bridge, nil
}

func (c *Controller) AddBridge(bridge *huego.Bridge) error {
	c.bridges[bridge.ID] = bridge
	dbo := Bridge{
		Url:      bridge.Host,
		Username: bridge.User,
	}
	result, err := c.db.Model(&dbo).Insert()
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (c *Controller) SetStateOnSpecificLight(bridgeId string, lightId int, state huego.State) error {
	return nil
}
