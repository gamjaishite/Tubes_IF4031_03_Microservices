import { db } from "../../configs/drizzle";
import { ServiceResponse } from "../../types/common";
import { BookingStatus, booking } from "../../models/booking";
import { GetBookingServicePayload } from "../../dto/booking/getBooking.dto";
import { CancelBookingData } from "../../types/booking";
import { TicketServiceResponse } from "../../types/common";
import { and, eq } from "drizzle-orm";

const cancelBookingService = async ({
	userId,
	bookingId,
	jwt,

}: GetBookingServicePayload): Promise<ServiceResponse> => {

    const bookingData = await db.query.booking.findFirst({
		where: and(
			eq(booking.userId, userId),
			eq(booking.id, bookingId)
		)
	});

    if (!bookingData) {
		return {
			code: 404,
			message: "Booking history not found",
		};
	}

    if (bookingData.status !== BookingStatus.WAITING_FOR_PAYMENT) {
        return {
            code: 403,
            message: "Cannot cancel ongoing payment"
        }
    }

    await fetch(`${process.env.TICKET_SERVICE_BASE_URL}/tikets/${bookingData.ticketId}/status/cancelled`, {
		method: 'PATCH',
		headers: {
			Authorization: `Bearer ${jwt}`
		},
		credentials: 'include'
	});

    const newStatus = BookingStatus.FAILED;

	await db
        .update(booking)
        .set({status: newStatus})
        .where(and(eq(booking.id, bookingId), eq(booking.userId, userId)));

    const cancelData: CancelBookingData = {
        newStatus
    }

	return {
		code: 200,
		message: "Success",
        data: cancelData
	}
};

export default cancelBookingService;
