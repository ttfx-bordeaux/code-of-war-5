package main

import (
	"image/color"
	"log"

	"github.com/engoengine/ecs"
	"github.com/engoengine/engo"
)

//GameWorld world
type GameWorld struct{}

//Preload preload
func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static", false)

	log.Println("Preloaded")
}

//Setup setup
func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})

	createGround(w, 3, 5, 0, "grass-600-600.png")
	createGround(w, 3, 5, 380, "stone-600-400.png")
	createGround(w, 3, 5, 780, "water-600-600.png")

	createTower(w, 1, 3, 0, "tour1-600-600.png")
}

func createTower(w *ecs.World, abs, ord float32, padding float32, imgName string) {
	tower := createEntityTile(imgName, engo.Point{X: abs*120 + padding, Y: ord * 120})
	err := w.AddEntity(tower)
	if err != nil {
		log.Println(err)
	}
}

func createGround(w *ecs.World, width, length int, padding float32, imgName string) {
	for j := 0; j < length; j++ {
		for i := 0; i < width; i++ {
			grass := createEntityTile(imgName, engo.Point{X: float32(i)*120 + padding, Y: float32(j) * 120})
			err := w.AddEntity(grass)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createEntityTile(imgName string, point engo.Point) *ecs.Entity {
	// Create an entity part of the Render
	entityTile := ecs.NewEntity([]string{"RenderSystem"})
	// Retrieve a texture
	texture := engo.Files.Image(imgName)
	// renvoie nill si image pas chargée
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}

	render := engo.NewRenderComponent(texture, engo.Point{0.2, 0.2}, "tile")

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{point, width, height}

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

func main() {
	opts := engo.RunOptions{
		Title:  "Code of War : Enlarge your tower",
		Width:  1024,
		Height: 640,
	}
	engo.Run(opts, &GameWorld{})
}
