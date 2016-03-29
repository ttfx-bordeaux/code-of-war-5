package main

import (
	"image/color"
	"log"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

/*
http://www.gamedev.net/page/resources/_/technical/game-programming/understanding-component-entity-systems-r3013
Entity Component System:
* Entity: The entity is a general purpose object. Usually, it only consists of a unique id.
    They "tag every coarse gameobject as a separate item".
    Implementations typically use a plain integer for this.
* Component: the raw data for one aspect of the object, and how it interacts with the world.
    "Labels the Entity as possessing this particular aspect".
    Implementations typically use Structs, Classes, or Associative Arrays.
* System: "Each System runs continuously (as though each System had its own private thread)
    and performs global actions on every Entity that possesses a Component of the same aspect as that System."
*/

//GameWorld world
type GameWorld struct{}

//Preload preload
func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engi.Files.AddFromDir("static", false)

	log.Println("Preloaded")
}

//Setup setup
func (game *GameWorld) Setup(w *ecs.World) {
	engi.SetBg(color.White)

	w.AddSystem(&engi.RenderSystem{})

	grass := createEntityTile("grass-600-600.png", engi.Point{0, 0})
	grassBold := createEntityTile("grass-bold-500-300.png", engi.Point{120, 0})
	stone := createEntityTile("stone-600-400.png", engi.Point{0, 120})
	water := createEntityTile("water-600-600.png", engi.Point{120, 120})

	err := w.AddEntity(grass)
	if err != nil {
		log.Println(err)
	}
	w.AddEntity(stone)
	w.AddEntity(grassBold)
	w.AddEntity(water)
}

func createEntityTile(imgName string, point engi.Point) *ecs.Entity {
	// Create an entity part of the Render
	guy := ecs.NewEntity([]string{"RenderSystem"})
	// Retrieve a texture
	texture := engi.Files.Image(imgName)
	// renvoie nill si image pas chargée

	// Create RenderComponent... Set scale to 0.2x, give lable "guy"
	render := engi.NewRenderComponent(texture, engi.Point{0.2, 0.2}, "guy")

	log.Println(texture)

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y

	space := &engi.SpaceComponent{point, width, height}

	guy.AddComponent(render)
	guy.AddComponent(space)

	return guy
}

//Hide hide
func (*GameWorld) Hide() {}

//Show show
func (*GameWorld) Show() {}

//Type type
func (*GameWorld) Type() string { return "GameWorld" }

// see https://github.com/paked/engi
func main() {
	opts := engi.RunOptions{
		Title:  "Code of War : Enlarge your tower",
		Width:  1024,
		Height: 640,
	}
	engi.Run(opts, &GameWorld{})
}
