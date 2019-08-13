/area/var/music

/area/Entered(mob/m)
    if (!ismob(m))
        return
    m << desc
    m << sound(music, 1, channel = 1)

/area/outside/desc = "Nice and jazzy, here..."
/area/outside/music = 'jazzy.mid'
/area/cave/desc = "Watch out for the giant rat!"
/area/cave/music = 'cavern.mid'

/mob/player/Bump(atom/obstacle)
    src << "You bump into [obstacle]."
    src << 'ouch.wav'
