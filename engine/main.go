package main

import (
	"image/color"
	"log"
	"time"

	"github.com/ttfx-bordeaux/code-of-war-5/engine/game"

	"engo.io/ecs"
	"engo.io/engo"
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

func toPoint(p game.PositionComponent, loop int) engo.Point {
	padding := loop * (nbTilesAbs*TileWidth + padRight)
	return engo.Point{X: float32(p.Abs*TileWidth + padding), Y: float32(p.Ord * TileHeight)}
}

type Player struct {
	Id       string
	Ground   []*ecs.BasicEntity
	towers   []*ecs.BasicEntity
	chickens []*ecs.BasicEntity
}

type towerEntity struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

type chickenEntity struct {
	ecs.BasicEntity
	engo.AnimationComponent
	engo.RenderComponent
	engo.SpaceComponent
	game.PositionComponent
}

type tileEntity struct {
	ecs.BasicEntity
	engo.RenderComponent
	engo.SpaceComponent
}

//GameWorld world
type DefaultScene struct{}

//Preload preload
func (world *DefaultScene) Preload() {
	// Load all files from the data directory. Do not do it recursively.
	engo.Files.AddFromDir("static", false)

	game.Preload()

	log.Println("Preloaded")
}

//Setup setup
func (world *DefaultScene) Setup(w *ecs.World) {
	engo.SetBackground(backgroundColor)

	renderSystem := &engo.RenderSystem{}
	animationSystem := &engo.AnimationSystem{}
	controlSystem := &game.ControlSystem{}

	w.AddSystem(renderSystem)
	w.AddSystem(animationSystem)
	w.AddSystem(controlSystem)

	idGround := 0
	var chicken *chickenEntity
	for _, p := range players {
		log.Println(p)
		log.Println(idGround)
		createGround(w, idGround, ImgGroundName, renderSystem)
		createTower(w, game.PositionComponent{1, 3}, idGround, ImgTowerName, renderSystem)
		chicken = createChicken(w, game.PositionComponent{0, 0}, idGround, renderSystem, animationSystem, controlSystem)
		idGround++
	}

	go loopGame(chicken, 1)
}

func loopGame(chicken *chickenEntity, idGround int) {
	// loop game
	//var time *engo.Clock
	ticker := time.NewTicker(time.Duration(int(time.Second)))
Outer:
	for {
		select {
		case <-ticker.C:
			cf := &chicken.PositionComponent
			space := &chicken.SpaceComponent

			if cf.Ord < nbTilesOrd - 1 {
				cf.Ord++
				space.Position = toPoint(*cf, idGround)
				log.Printf("move")
			} else {
				break Outer
			}
			// call http
			//if close {
			//	break Outer
			// }
			// if !headless && window.ShouldClose() {
			// closeEvent()
			// }
		}
	}
	log.Println("stop")
	ticker.Stop()
}

func createTower(w *ecs.World, p game.PositionComponent, idGround int, imgName string, renderSystem *engo.RenderSystem) {
	tower := createEntityTower(imgName, toPoint(p, idGround))
	renderSystem.Add(&tower.BasicEntity, &tower.RenderComponent, &tower.SpaceComponent)
}

func createGround(w *ecs.World, idGround int, imgName string, renderSystem *engo.RenderSystem) {
	padding := idGround * (nbTilesAbs*TileWidth + padRight)
	for j := 0; j < nbTilesOrd; j++ {
		for i := 0; i < nbTilesAbs; i++ {
			tile := createEntityTile(imgName, engo.Point{X: float32(i*TileWidth + padding), Y: float32(j * TileHeight)})
			renderSystem.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
		}
	}
}

func createChicken(w *ecs.World, p game.PositionComponent, idGround int,
	renderSystem *engo.RenderSystem, animationSystem *engo.AnimationSystem,
	controlSystem *game.ControlSystem) *chickenEntity {
	spriteSheet := engo.NewSpritesheetFromFile("chicken.png", 150, 150)
	entity := createEntityChicken(toPoint(p, idGround), spriteSheet)

	// Add our hero to the appropriate systems
	renderSystem.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)
	animationSystem.Add(&entity.BasicEntity, &entity.AnimationComponent, &entity.RenderComponent)
	controlSystem.Add(&entity.BasicEntity, &entity.AnimationComponent)

	return entity;
}

func createEntityChicken(point engo.Point, spriteSheet *engo.Spritesheet) *chickenEntity {
	entity := &chickenEntity{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: 150, Height: 150}
	entity.RenderComponent = engo.NewRenderComponent(spriteSheet.Cell(0), engo.Point{X: 1, Y: 1})
	entity.AnimationComponent = engo.NewAnimationComponent(spriteSheet.Drawables(), game.AnimationRate)
	entity.AnimationComponent.AddAnimations(game.Actions)
	entity.AnimationComponent.AddDefaultAnimation(game.StopAction)

	return entity
}

func createEntityTile(imgName string, point engo.Point) *tileEntity {
	entity := &tileEntity{BasicEntity: ecs.NewBasic()}

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}

	entity.RenderComponent = engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2})
	entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: texture.Width(), Height: texture.Height()}

	return entity
}

func createEntityTower(imgName string, point engo.Point) *towerEntity {
	entity := &towerEntity{BasicEntity: ecs.NewBasic()}

	texture := engo.Files.Image(imgName)
	if texture == nil {
		log.Fatalf("image %s not loaded\n", imgName)
	}
	entity.SpaceComponent = engo.SpaceComponent{Position: point, Width: texture.Width(), Height: texture.Height()}
	entity.RenderComponent = engo.NewRenderComponent(texture, engo.Point{X: 0.2, Y: 0.2})

	return entity
}

//Hide hide
func (*DefaultScene) Hide() {}

//Show show
func (*DefaultScene) Show() {}

//Type type
func (*DefaultScene) Type() string { return "GameWorld" }

func main() {
	// see level.go
	// tick, fps
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
	engo.Run(opts, &DefaultScene{})
}
