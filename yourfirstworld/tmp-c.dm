
/mob/player/verb/look()
    src << "You see..."

    for(var/atom/movable/o in oview())
        src << "[o].  [o.desc]"
