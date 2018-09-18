package main

func main() {

	/*
		// this is old prototype code. ignore this for now.
		// FIXME: reimplement this
		ctx, cancel := context.WithCancel(context.Background())
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)

		room := rooms.PlayRoom{}
		zombie := &zombies.Easy{}

		room.AddZombie(zombie)

		room.Start(ctx)

		for {
			select {
			case event := <-room.EventStream():
				log.Println(event.String())
			case <-stop:
				cancel()
				return
			}
		}
	*/

}
