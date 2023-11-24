// app.d.ts
/// <reference types="lucia" />
declare namespace Lucia {
	type Auth = import("./src/configs/lucia.js").Auth;
	type DatabaseUserAttributes = {
		email: string;
	};
	type DatabaseSessionAttributes = {};
}
