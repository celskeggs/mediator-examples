
turf
    verb
        greet_contents()
            set src in world
            src << "you should see this if you're standing here"
            list(src) << "you should NOT see this if you're standing here"

mob
	player
		verb
            say2(msg as text)
                for(var/v in view(src))
                    v << "[src] says, \"[msg]\" via [v]"
