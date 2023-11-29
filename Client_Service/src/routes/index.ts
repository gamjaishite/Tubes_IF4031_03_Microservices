import Elysia from "elysia";
import { authRouteV1 } from "./auth";
import { profileRouteV1 } from "./profile";
import { bookingRouteV1 } from "./booking";
import cors from "@elysiajs/cors";

const Routes = (app: Elysia) => {
	app
		.use(
			cors({
				origin: [/http:\/\/localhost:[0-9]{4}$/],
				credentials: true,
				methods: ["GET", "PUT", "POST", "DELETE", "PATCH"],
				allowedHeaders: ["Content-Type", "Authorization"],
			}),
		)
		.group("/api", (app) =>
			app.group("/v1", (app) => app.use(profileRouteV1).use(authRouteV1).use(bookingRouteV1)),
		);
};

export default Routes;6
