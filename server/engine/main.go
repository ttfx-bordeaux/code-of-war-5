package main

import (
	"image/color"
	"log"

	"github.com/paked/engi"
	"github.com/paked/engi/ecs"
)

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

	createGround(w, 3, 5, 0, "grass-600-600.png")
	createGround(w, 3, 5, 380, "stone-600-400.png")
	createGround(w, 3, 5, 780, "water-600-600.png")
}

func createGround(w *ecs.World, width, length int, padding float32, imgName string) {
	for j := 0; j < length; j++ {
		for i := 0; i < width; i++ {
			grass := createEntityTile(imgName, engi.Point{X: float32(i)*120 + padding, Y: float32(j) * 120})
			err := w.AddEntity(grass)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createEntityTile(imgName string, point engi.Point) *ecs.Entity {
	// Create an entity part of the Render
	entityTile := ecs.NewEntity([]string{"RenderSystem"})
	// Retrieve a texture
	texture := engi.Files.Image(imgName)
	// renvoie nill si image pas chargée
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}

	render := engi.NewRenderComponent(texture, engi.Point{0.2, 0.2}, "tile")

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engi.SpaceComponent{point, width, height}

	entityTile.AddComponent(render)
	entityTile.AddComponent(space)

	return entityTile
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
