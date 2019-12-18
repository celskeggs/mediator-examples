obj
	cheese
		desc = "It is quite smelly."

	scroll
		desc = "It looks to be rather old."

/mob/player/desc = "A handsome and dashing rogue."

/mob/player/verb/look()
    src << "You see..."

    for(var/atom/movable/o in oview())
        src << "[o].  [o.desc]"

/mob/rat
	desc = "It's quite large."
