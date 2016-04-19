package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"github.com/ttfx-bordeaux/code-of-war-5/engine/game"
)

var (
	screenWidth                 = 1024
	screenHeight                = 640
	TileWidth       int         = 120
	TileHeight      int         = 120
	nbTilesAbs      int         = 3
	nbTilesOrd      int         = 5
	ImgGroundName   string      = "grass-600-600.png"
	ImgTowerName    string      = "tour1-600-600.png"
	backgroundColor color.Color = color.White
	padRight        int
	players         map[string]*Player
)

type GamePosition struct {
	Abs int
	Ord int
}

func (p GamePosition) toPoint(loop int) engo.Point {
	padding := loop * (nbTilesAbs*TileWidth + padRight)
	return engo.Point{X: float32(p.Abs*TileWidth + padding), Y: float32(p.Ord * TileHeight)}
}

type Player struct {
	Id       string
	Ground   []*ecs.Entity
	towers   []*ecs.Entity
	chickens []*ecs.Entity
}

//GameWorld world
type GameWorld struct{}

//Preload preload
func (game *GameWorld) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static", false)

	system.Preload()

	log.Println("Preloaded")
}

//Setup setup
func (game *GameWorld) Setup(w *ecs.World) {
	engo.SetBackground(backgroundColor)

	w.AddSystem(&engo.RenderSystem{})
	w.AddSystem(&engo.AnimationSystem{})
	w.AddSystem(&system.ControlSystem{})

	loop := 0
	for _, p := range players {
		log.Println(p)
		log.Println(loop)
		createGround(w, loop, ImgGroundName)
		createTower(w, GamePosition{1, 3}, loop, ImgTowerName)
		createChicken(w, GamePosition{0, 0}, loop)
		loop++
	}
	// entities
	//createGround(w, 0, ImgGroundName)
	//createGround(w, 3, 5, 380, "stone-600-400.png")
	//createGround(w, 3, 5, 780, "water-600-600.png")
}

func createTower(w *ecs.World, p GamePosition, loop int, imgName string) {
	tower := createEntityTower(imgName, p.toPoint(loop))
	err := w.AddEntity(tower)
	if err != nil {
		log.Println(err)
	}
}

func createGround(w *ecs.World, loop int, imgName string) {
	padding := loop * (nbTilesAbs*TileWidth + padRight)
	for j := 0; j < nbTilesOrd; j++ {
		for i := 0; i < nbTilesAbs; i++ {
			grass := createEntityTile(imgName, engo.Point{X: float32(i*TileWidth + padding), Y: float32(j * TileHeight)})
			err := w.AddEntity(grass)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func createChicken(w *ecs.World, p GamePosition, loop int) {
	spriteSheet := engo.NewSpritesheetFromFile("chicken.png", 150, 150)
	err := w.AddEntity(createEntityChicken(p.toPoint(loop), spriteSheet, system.StopAction))
	if err != nil {
		log.Println(err)
	}
}

func createEntityChicken(point engo.Point, spriteSheet *engo.Spritesheet, action *engo.AnimationAction) *ecs.Entity {
	entity := ecs.NewEntity("AnimationSystem", "RenderSystem", "ControlSystem")

	render := engo.NewRenderComponent(spriteSheet.Cell(action.Frames[0]), engo.Point{X: 1, Y: 1}, "chicken")
	width := 150 * render.Scale().X
	height := 150 * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}
	animation := engo.NewAnimationComponent(spriteSheet.Drawables(), system.AnimationRate)
	animation.AddAnimationActions(system.Actions)
	animation.SelectAnimationByAction(action)

	entity.AddComponent(render)
	entity.AddComponent(space)
	entity.AddComponent(animation)

	return entity
}

func createEntityTile(imgName string, point engo.Point) *ecs.Entity {
	entity := ecs.NewEntity("RenderSystem")

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}
	render := engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2}, "tile")
	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}

	entity.AddComponent(render)
	entity.AddComponent(space)

	return entity
}

func createEntityTower(imgName string, point engo.Point) *ecs.Entity {
	entity := ecs.NewEntity("RenderSystem")

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}
	render := engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2}, "tile")
	width := texture.Width() * render.Scale().X
	height := texture.Height() * render.Scale().Y
	space := &engo.SpaceComponent{Position: point, Width: width, Height: height}

	entity.AddComponent(render)
	entity.AddComponent(space)

	return entity
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
	players = make(map[string]*Player, 2)
	players["1"] = &Player{Id: "1"}
	players["2"] = &Player{Id: "2"}
	sizeWidth := TileWidth * nbTilesAbs * len(players)
	if sizeWidth > screenWidth {
		padRight = 10
	} else {
		padRight = (screenWidth - sizeWidth) / len(players)
	}
	opts := engo.RunOptions{
		Title:      "Code of War : Enlarge your tower",
		Width:      screenWidth,
		Height:     screenHeight,
		Fullscreen: false,
	}
	engo.OverrideCloseAction()
	engo.Run(opts, &GameWorld{})
}
