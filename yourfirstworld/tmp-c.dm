obj
	verb
		get()
/*
			This makes it so that you have to be standing on top
			of the obj to pick it up. The oview() procedure is
			similar to view() in that it retrieves a list of
			everything in sight, but it doesn't include the center.
*/
			set src in oview(0)

			usr << "You get [src]."

			// Move the obj into the player's contents.
			Move(usr)

		drop()
			usr << "You drop [src]."

			Move(usr.loc)

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
