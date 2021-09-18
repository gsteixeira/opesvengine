package opesvengine

import (
    "image"
    "github.com/hajimehoshi/ebiten/v2"
)

// Sprite objects.
type Sprite struct {
    // shape
    Width int
    Heigth int
    scale_x float64
    scale_y float64
    kind string
    Parent *Character
    payload int // aux var used by shoot damage
    max_life int
    // movement
    pos_x float64
    pos_y float64
    speed_x float64
    speed_y float64
    // animation
    count int
    direction int
    Frame_offset_x int
    Frame_offset_y int
    frame_width int
    frame_heigth int
    frame_number int
    // physics
    grip_x float64
    grip_y float64
    on_the_ground bool
    solid bool
    Static bool
    gravity float64
    // ebiten
    eimage *ebiten.Image
    sub_img image.Image
}

// this struct keeps info about the sprite boundaries,
// so we dont have to calculate it again
type SpriteRect struct {
        left float64
        right float64
        top float64
        bottom float64
    }

// Sprite constructor method
func NewSprite(x float64, y float64, imgr image.Image, is_static bool) Sprite {
    s:= Sprite{
        // shape
        Width: FRAME_WIDTH,
        Heigth: FRAME_HEIGTH,
        scale_x: 1,
        scale_y: 1,
        kind: "default",
        Parent: nil,
        payload: 0,
        // movement
        pos_x: x,
        pos_y: y,
        speed_x: 0,
        speed_y: 0,
        // animation
        count: 0,
        direction: 0, // 1 = right, -1 left
        Frame_offset_x: 0,
        Frame_offset_y: 0,
        frame_width: FRAME_WIDTH,
        frame_heigth: FRAME_HEIGTH,
        frame_number: 0,
        // physics
        grip_x: 0,
        grip_y: 0,
        on_the_ground: false,
        solid: true,
        Static: is_static,
        max_life: -1,
        gravity: GRAVITY,
        // ebiten image
        eimage: ebiten.NewImageFromImage(imgr),
    }
    return s
}

// Set sprite frame (sub image and tile size) parameters
func (s *Sprite) Set_frames(frame_width int, frame_heigth int,
                            frame_number int) {
    s.frame_width = frame_width
    s.frame_heigth = frame_heigth
    s.frame_number = frame_number
}

// Set sprite dimensions
func (s *Sprite) Set_dimensions(width, heigth int) {
    s.Width = width
    s.Heigth = heigth
    
    s.Frame_offset_x = (s.frame_width - width) / 2
    s.Frame_offset_y = (s.frame_heigth - heigth) / 2
}

// Get position on map Y of top of sprite
func (s *Sprite) rect_top() float64 {
    return s.pos_y
}

// Get position on map Y at left of sprite
func (s *Sprite) rect_left() float64 {
    return s.pos_x
}

// Get position on map Y at right of sprite
func (s *Sprite) rect_right() float64 {
    return s.pos_x + float64(s.Width)
}

// Get position on map Y at botton of sprite
func (s *Sprite) rect_bottom() float64 {
    return s.pos_y + float64(s.Heigth)
}

// Get position of the center of the sprite
func (s *Sprite) pos_center() (float64, float64) {
    x := s.pos_x + (float64(s.Width) / 2)
    y := s.pos_y + (float64(s.Heigth) / 2)
    return x, y
}

// get sprite rectangle in a struct for eficiency
// for future movements, pass a speed > 0
func (s *Sprite) get_rect(speed_x float64, speed_y float64) SpriteRect {
    var rect = SpriteRect{
                s.rect_left() + speed_x,
                s.rect_right() + speed_x,
                s.rect_top() + speed_y,
                s.rect_bottom() + speed_y,
                }
    return rect
}

// Update animation parameters
// This changes the offset of the image map we use to take tiles from.
func (s *Sprite) setup_animation() bool {
    const RESTING = 5
    const RUNNING = 8
    const JUMPING = 4
    if !s.Static {
        if s.speed_y != 0 {
            s.Frame_offset_y = s.frame_heigth * 2
            s.frame_number = JUMPING
        } else if s.speed_x != 0 {
            s.Frame_offset_y = s.frame_heigth
            s.frame_number = RUNNING                  
        } else {
            s.Frame_offset_y = 0
            s.frame_number = RESTING
        }
        return true
    }
    // TODO flip image horizontally when moving left
    // find a performant way to do it. (probably is cheapest to have an image for that)
    // if s.direction < 0 {
    //      flip..
    // }
    return false
}

// Draw the sprite on the screen
func (s *Sprite) draw(screen *ebiten.Image) {
    s.setup_animation()
    var sx, sy, i int
    opp := &ebiten.DrawImageOptions{}
    opp.GeoM.Translate(s.pos_x, s.pos_y)
    opp.GeoM.Scale(s.scale_x, s.scale_y)
    // to animate or not to animate, that's the question.
    if s.Static {
        i = 1
    } else {
        i = (s.count / 5) % s.frame_number
    }
    sx, sy = s.Frame_offset_x + i * s.frame_width, s.Frame_offset_y
    img_rect := image.Rect(sx, sy, sx + s.frame_width, sy + s.frame_heigth)
    sub_img := s.eimage.SubImage(img_rect)
    screen.DrawImage(sub_img.(*ebiten.Image), opp)
}

// get X, Y position based on map instead of screen
func (s *Sprite) Get_position_onmap(game *Universe) (float64, float64) {
    return s.pos_x + float64(game.MapOffset_x), s.pos_y + float64(game.MapOffset_y)
}
