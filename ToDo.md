## Stuff To DO 


### Fixes 

1. If World and Viewport is same AVOID Viewport calculations (DONE)
2. Allow Static Viewport (when only screenshake is needed). (DONE)
3. Separate Constructor for Static Viewport (World padding auto? for screenshake) (DONE same constructor)
4. ViewportMargin > will maintain a min margin between viewport and world (to make sure padding area of world (to account for shake) is never visible in Viewport) (DONE)
5. Make ScreenOptions (LEFT), Viewport InitialPosition (DONE)
6. Remove Screen Offset calculation if possible. (DONE)
7. Privatize some fields. LEFT


FixedViewport = Viewport cannot move (viewport can be nil)
AutoPadding = Viewport and WorldView same dimensions  
AutoPadding = autoOffset (0,0) to (-padding, -padding)
FixedCamera = Camera Cannot move (camera can be nil)