package mqinterface

import (
	"github.com/streadway/amqp"
)

// Channel opens a unique, concurrent server channel to process the bulk of AMQP
// messages. Any error from methods on this receiver will render the receiver
// invalid and a new Channel should be opened.
type Channel interface {
	/*
		Confirm puts this channel into confirm mode so that the client can ensure all
		publishing's have successfully been received by the server. After entering this
		mode, the server will send a basic.ack or basic.nack message with the deliver
		tag set to a 1 based incrementing index corresponding to every publishing
		received after the this method returns.

		Add a listener to Channel.NotifyPublish to respond to the Confirmations. If
		Channel.NotifyPublish is not called, the Confirmations will be silently
		ignored.

		The order of acknowledgments is not bound to the order of deliveries.

		Ack and Nack confirmations will arrive at some point in the future.

		Unroutable mandatory or immediate messages are acknowledged immediately after
		any Channel.NotifyReturn listeners have been notified. Other messages are
		acknowledged when all queues that should have the message routed to them have
		either have received acknowledgment of delivery or have enqueued the message,
		persisting the message if necessary.

		When noWait is true, the client will not wait for a response. A channel
		exception could occur if the server does not support this method.

	*/
	Confirm(noWait bool) error

	/*
	   NotifyPublish registers a listener for reliable publishing. Receives from this
	   chan for every publish after Channel.Confirm will be in order starting with
	   DeliveryTag 1.

	   There will be one and only one Confirmation Publishing starting with the
	   delivery tag of 1 and progressing sequentially until the total number of
	   publishing's have been seen by the server.

	   Acknowledgments will be received in the order of delivery from the
	   NotifyPublish channels even if the server acknowledges them out of order.

	   The listener chan will be closed when the Channel is closed.

	   The capacity of the chan Confirmation must be at least as large as the
	   number of outstanding publishing's. Not having enough buffered chans will
	   create a deadlock if you attempt to perform other operations on the Connection
	   or Channel while confirms are in-flight.

	   It's advisable to wait for all Confirmations to arrive before calling
	   Channel.Close() or Connection.Close().

	*/
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation

	/*
	   Publish sends a Publishing from the client to an exchange on the server.

	   When you want a single message to be delivered to a single queue, you can
	   publish to the default exchange with the routingKey of the queue name. This is
	   because every declared queue gets an implicit route to the default exchange.

	   Since publishing's are asynchronous, any undeliverable message will get returned
	   by the server. Add a listener with Channel.NotifyReturn to handle any
	   undeliverable message when calling publish with either the mandatory or
	   immediate parameters as true.

	   publishing's can be undeliverable when the mandatory flag is true and no queue is
	   bound that matches the routing key, or when the immediate flag is true and no
	   consumer on the matched queue is ready to accept the delivery.

	   This can return an error when the channel, connection or socket is closed. The
	   error or lack of an error does not indicate whether the server has received this
	   publishing.

	   It is possible for publishing to not reach the broker if the underlying socket
	   is shutdown without pending publishing packets being flushed from the kernel
	   buffers. The easy way of making it probable that all publishing's reach the
	   server is to always call Connection.Close before terminating your publishing
	   application. The way to ensure that all publishing's reach the server is to add
	   a listener to Channel.NotifyPublish and put the channel in confirm mode with
	   Channel.Confirm. Publishing delivery tags and their corresponding
	   confirmations start at 1. Exit when all publishing's are confirmed.

	   When Publish does not return an error and the channel is in confirm mode, the
	   internal counter for DeliveryTags with the first confirmation starting at 1.

	*/
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error

	/*
	   Close initiate a clean channel closure by sending a close message with the error
	   code set to '200'.

	   It is safe to call this method multiple times.

	*/
	Close() error

	/*
		ExchangeDeclare declares an exchange on the server. If the exchange does not
		already exist, the server will create it.  If the exchange exists, the server
		verifies that it is of the provided type, durability and auto-delete flags.

		Errors returned from this method will close the channel.

		Exchange names starting with "amq." are reserved for pre-declared and
		standardized exchanges. The client MAY declare an exchange starting with
		"amq." if the passive option is set, or the exchange already exists.  Names can
		consist of a non-empty sequence of letters, digits, hyphen, underscore,
		period, or colon.

		Each exchange belongs to one of a set of exchange kinds/types implemented by
		the server. The exchange types define the functionality of the exchange - i.e.
		how messages are routed through it. Once an exchange is declared, its type
		cannot be changed.  The common types are "direct", "fanout", "topic" and
		"headers".

		Durable and Non-Auto-Deleted exchanges will survive server restarts and remain
		declared when there are no remaining bindings.  This is the best lifetime for
		long-lived exchange configurations like stable routes and default exchanges.

		Non-Durable and Auto-Deleted exchanges will be deleted when there are no
		remaining bindings and not restored on server restart.  This lifetime is
		useful for temporary topologies that should not pollute the virtual host on
		failure or after the consumers have completed.

		Non-Durable and Non-Auto-deleted exchanges will remain as long as the server is
		running including when there are no remaining bindings.  This is useful for
		temporary topologies that may have long delays between bindings.

		Durable and Auto-Deleted exchanges will survive server restarts and will be
		removed before and after server restarts when there are no remaining bindings.
		These exchanges are useful for robust temporary topologies or when you require
		binding durable queues to auto-deleted exchanges.

		Note: RabbitMQ declares the default exchange types like 'amq.fanout' as
		durable, so queues that bind to these pre-declared exchanges must also be
		durable.

		Exchanges declared as `internal` do not accept accept publishings. Internal
		exchanges are useful when you wish to implement inter-exchange topologies
		that should not be exposed to users of the broker.

		When noWait is true, declare without waiting for a confirmation from the server.
		The channel may be closed as a result of an error.  Add a NotifyClose listener
		to respond to any exceptions.

		Optional amqp.Table of arguments that are specific to the server's implementation of
		the exchange can be sent for exchange types that require extra parameters.
	*/
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error

	/*
		QueueDeclare declares a queue to hold messages and deliver to consumers.
		Declaring creates a queue if it doesn't already exist, or ensures that an
		existing queue matches the same parameters.

		Every queue declared gets a default binding to the empty exchange "" which has
		the type "direct" with the routing key matching the queue's name.  With this
		default binding, it is possible to publish messages that route directly to
		this queue by publishing to "" with the routing key of the queue name.

		  QueueDeclare("alerts", true, false, false, false, nil)
		  Publish("", "alerts", false, false, Publishing{Body: []byte("...")})

		  Delivery       Exchange  Key       Queue
		  -----------------------------------------------
		  key: alerts -> ""     -> alerts -> alerts

		The queue name may be empty, in which case the server will generate a unique name
		which will be returned in the Name field of Queue struct.

		Durable and Non-Auto-Deleted queues will survive server restarts and remain
		when there are no remaining consumers or bindings.  Persistent publishings will
		be restored in this queue on server restart.  These queues are only able to be
		bound to durable exchanges.

		Non-Durable and Auto-Deleted queues will not be redeclared on server restart
		and will be deleted by the server after a short time when the last consumer is
		canceled or the last consumer's channel is closed.  Queues with this lifetime
		can also be deleted normally with QueueDelete.  These durable queues can only
		be bound to non-durable exchanges.

		Non-Durable and Non-Auto-Deleted queues will remain declared as long as the
		server is running regardless of how many consumers.  This lifetime is useful
		for temporary topologies that may have long delays between consumer activity.
		These queues can only be bound to non-durable exchanges.

		Durable and Auto-Deleted queues will be restored on server restart, but without
		active consumers will not survive and be removed.  This Lifetime is unlikely
		to be useful.

		Exclusive queues are only accessible by the connection that declares them and
		will be deleted when the connection closes.  Channels on other connections
		will receive an error when attempting  to declare, bind, consume, purge or
		delete a queue with the same name.

		When noWait is true, the queue will assume to be declared on the server.  A
		channel exception will arrive if the conditions are met for existing queues
		or attempting to modify an existing queue from a different connection.

		When the error return value is not nil, you can assume the queue could not be
		declared with these parameters, and the channel will be closed.

	*/
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)

	/*
		Qos controls how many messages or how many bytes the server will try to keep on
		the network for consumers before receiving delivery acks.  The intent of Qos is
		to make sure the network buffers stay full between the server and client.

		With a prefetch count greater than zero, the server will deliver that many
		messages to consumers before acknowledgments are received.  The server ignores
		this option when consumers are started with noAck because no acknowledgments
		are expected or sent.

		With a prefetch size greater than zero, the server will try to keep at least
		that many bytes of deliveries flushed to the network before receiving
		acknowledgments from the consumers.  This option is ignored when consumers are
		started with noAck.

		When global is true, these Qos settings apply to all existing and future
		consumers on all channels on the same connection.  When false, the Channel.Qos
		settings will apply to all existing and future consumers on this channel.
		RabbitMQ does not implement the global flag.

		To get round-robin behavior between consumers consuming from the same queue on
		different connections, set the prefetch count to 1, and the next available
		message on the server will be delivered to the next available consumer.

		If your consumer work time is reasonably consistent and not much greater
		than two times your network round trip time, you will see significant
		throughput improvements starting with a prefetch count of 2 or slightly
		greater as described by benchmarks on RabbitMQ.

		http://www.rabbitmq.com/blog/2012/04/25/rabbitmq-performance-measurements-part-2/
	*/
	Qos(prefetchCount, prefetchSize int, global bool) error

	/*
		QueueBind binds an exchange to a queue so that publishings to the exchange will
		be routed to the queue when the publishing routing key matches the binding
		routing key.

		  QueueBind("pagers", "alert", "log", false, nil)
		  QueueBind("emails", "info", "log", false, nil)

		  Delivery       Exchange  Key       Queue
		  -----------------------------------------------
		  key: alert --> log ----> alert --> pagers
		  key: info ---> log ----> info ---> emails
		  key: debug --> log       (none)    (dropped)

		If a binding with the same key and arguments already exists between the
		exchange and queue, the attempt to rebind will be ignored and the existing
		binding will be retained.

		In the case that multiple bindings may cause the message to be routed to the
		same queue, the server will only route the publishing once.  This is possible
		with topic exchanges.

		  QueueBind("pagers", "alert", "amq.topic", false, nil)
		  QueueBind("emails", "info", "amq.topic", false, nil)
		  QueueBind("emails", "#", "amq.topic", false, nil) // match everything

		  Delivery       Exchange        Key       Queue
		  -----------------------------------------------
		  key: alert --> amq.topic ----> alert --> pagers
		  key: info ---> amq.topic ----> # ------> emails
		                           \---> info ---/
		  key: debug --> amq.topic ----> # ------> emails

		It is only possible to bind a durable queue to a durable exchange regardless of
		whether the queue or exchange is auto-deleted.  Bindings between durable queues
		and exchanges will also be restored on server restart.

		If the binding could not complete, an error will be returned and the channel
		will be closed.

		When noWait is true and the queue could not be bound, the channel will be
		closed with an error.

	*/
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error

	/*
		Consume immediately starts delivering queued messages.

		Begin receiving on the returned chan Delivery before any other operation on the
		Connection or Channel.

		Continues deliveries to the returned chan Delivery until Channel.Cancel,
		Connection.Close, Channel.Close, or an AMQP exception occurs.  Consumers must
		range over the chan to ensure all deliveries are received.  Unreceived
		deliveries will block all methods on the same connection.

		All deliveries in AMQP must be acknowledged.  It is expected of the consumer to
		call Delivery.Ack after it has successfully processed the delivery.  If the
		consumer is cancelled or the channel or connection is closed any unacknowledged
		deliveries will be requeued at the end of the same queue.

		The consumer is identified by a string that is unique and scoped for all
		consumers on this channel.  If you wish to eventually cancel the consumer, use
		the same non-empty identifier in Channel.Cancel.  An empty string will cause
		the library to generate a unique identity.  The consumer identity will be
		included in every Delivery in the ConsumerTag field

		When autoAck (also known as noAck) is true, the server will acknowledge
		deliveries to this consumer prior to writing the delivery to the network.  When
		autoAck is true, the consumer should not call Delivery.Ack. Automatically
		acknowledging deliveries means that some deliveries may get lost if the
		consumer is unable to process them after the server delivers them.
		See http://www.rabbitmq.com/confirms.html for more details.

		When exclusive is true, the server will ensure that this is the sole consumer
		from this queue. When exclusive is false, the server will fairly distribute
		deliveries across multiple consumers.

		The noLocal flag is not supported by RabbitMQ.

		It's advisable to use separate connections for
		Channel.Publish and Channel.Consume so not to have TCP pushback on publishing
		affect the ability to consume messages, so this parameter is here mostly for
		completeness.

		When noWait is true, do not wait for the server to confirm the request and
		immediately begin deliveries.  If it is not possible to consume, a channel
		exception will be raised and the channel will be closed.

		Optional arguments can be provided that have specific semantics for the queue
		or server.

		When the channel or connection closes, all delivery chans will also close.

		Deliveries on the returned chan will be buffered indefinitely. To limit memory
		of this buffer, use the Channel.Qos method to limit the amount of
		unacknowledged/buffered deliveries the server will deliver on this Channel.

	*/
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)

	/*
		NotifyClose registers a listener for when the server sends a channel or
		connection exception in the form of a Connection.Close or Channel.Close method.
		Connection exceptions will be broadcast to all open channels and all channels
		will be closed, where channel exceptions will only be broadcast to listeners to
		this channel.

		The chan provided will be closed when the Channel is closed and on a
		graceful close, no error will be sent.

	*/
	NotifyClose(c chan *amqp.Error) chan *amqp.Error

	/*
		Cancel stops deliveries to the consumer chan established in Channel.Consume and
		identified by consumer.

		Only use this method to cleanly stop receiving deliveries from the server and
		cleanly shut down the consumer chan identified by this tag.  Using this method
		and waiting for remaining messages to flush from the consumer chan will ensure
		all messages received on the network will be delivered to the receiver of your
		consumer chan.

		Continue consuming from the chan Delivery provided by Channel.Consume until the
		chan closes.

		When noWait is true, do not wait for the server to acknowledge the cancel.
		Only use this when you are certain there are no deliveries in flight that
		require an acknowledgment, otherwise they will arrive and be dropped in the
		client without an ack, and will not be redelivered to other consumers.

	*/
	Cancel(consumer string, noWait bool) error
}
