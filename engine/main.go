package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
)

var (
	TileWidth float32
	TileHeight float32
	animationRate float32
	WalkAction *engo.AnimationAction
	StopAction *engo.AnimationAction
	DieAction  *engo.AnimationAction
	actions    []*engo.AnimationAction
)

//GameWorld world
type GameWorld struct{}

//Preload preload
func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static", false)

	//engo.Files.Add("static/chicken.png")
	// animation for chicken
	StopAction = &engo.AnimationAction{Name: "stop", Frames: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	WalkAction = &engo.AnimationAction{Name: "move", Frames: []int{11, 12, 13, 14, 15}}
	DieAction = &engo.AnimationAction{Name: "die", Frames: []int{28, 29, 30}}
	actions = []*engo.AnimationAction{DieAction, StopAction, WalkAction}
	
	log.Println("Preloaded")
}

//Setup setup
func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(color.White)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&ControlSystem{})
	
	// entities
	createGround(w, 3, 5, 0, "grass-600-600.png")
	//createGround(w, 3, 5, 380, "stone-600-400.png")
	//createGround(w, 3, 5, 780, "water-600-600.png")

	createTower(w, 1, 3, 0, "tour1-600-600.png")
	
	createChicken(w)
}

func createTower(w *ecs.World, abs, ord float32, padding float32, imgName string) {
	tower := createEntityTile(imgName, engo.Point{X: abs*TileWidth + padding, Y: ord * TileHeight})
	err := w.AddEntity(tower)
	if err != nil {
		log.Println(err)
	}
}

func createGround(w *ecs.World, width, length int, padding float32, imgName string) {
	for j := 0; j < length; j++ {
		for i := 0; i < width; i++ {
			grass := createEntityTile(imgName, engo.Point{X: float32(i)*TileWidth + padding, Y: float32(j) * TileHeight})
			err := w.AddEntity(grass)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createChicken(w *ecs.World) {
	spriteSheet := engo.NewSpritesheetFromFile("chicken.png", 150, 150)

	err := w.AddEntity(createEntity(&engo.Point{X:0, Y:0}, spriteSheet, StopAction))
	if err != nil {
		log.Println(err)
	}
}

func createEntity(point *engo.Point, spriteSheet *engo.Spritesheet, action *engo.AnimationAction) *ecs.Entity {
	entity := ecs.NewEntity("AnimationSystem", "RenderSystem", "ControlSystem")

	
	render := engo.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engo.Point{X:1, Y:1}, "chicken")
	width := 150 * render.Scale().X
	height := 150 * render.Scale().Y
	space := &engo.SpaceComponent{Position:*point, Width:width, Height:height}
	animation := engo.NewAnimationComponent(spriteSheet.Drawables(), animationRate)
	animation.AddAnimationActions(actions)
	animation.SelectAnimationByAction(action)
	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

func createEntityTile(imgName string, point engo.Point) *ecs.Entity {
	// Create an entity part of the Render
	entityTile := ecs.NewEntity("RenderSystem")
	// Retrieve a texture
	texture := engo.Files.Image(imgName)
	// renvoie nill si image pas chargée
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}

	render := engo.NewRenderComponent(texture, engo.Point{X:0.2, Y:0.2}, "tile")

	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{Position:point, Width:width, Height:height}

	entityTile.AddComponent(render)
	entityTile.AddComponent(space)

	return entityTile
}

type ControlSystem struct {
	ecs.LinearSystem
}

func (*ControlSystem) Type() string { return "ControlSystem" }
func (*ControlSystem) Pre()         {}
func (*ControlSystem) Post()        {}

func (c *ControlSystem) New(*ecs.World) {}

func (c *ControlSystem) UpdateEntity(entity *ecs.Entity, dt float32) {
	var a *engo.AnimationComponent

	if !entity.Component(&a) {
		return
	}

	if engo.Keys.Get(engo.ArrowRight).Down() {
		a.SelectAnimationByAction(WalkAction)
	} else if engo.Keys.Get(engo.Space).Down() {
		a.SelectAnimationByAction(DieAction)
	} else {
		a.SelectAnimationByAction(StopAction)
	}
}

//Hide hide
func (*GameWorld) Hide() {}

//Show show
func (*GameWorld) Show() {}

//Exit
func (*GameWorld) Exit() {
	log.Println("[GAME] Exit event called")
	//Here if you want you can prompt the user if they're sure they want to close
	engo.Exit()
}

//Type type
func (*GameWorld) Type() string { return "GameWorld" }

func main() {
	opts := engo.RunOptions{
		Title:  "Code of War : Enlarge your tower",
		Width:  1024,
		Height: 640,
	}
	TileHeight= 120
	TileWidth= 120
	animationRate=0.1
	engo.OverrideCloseAction()
	engo.Run(opts, &GameWorld{})
}
