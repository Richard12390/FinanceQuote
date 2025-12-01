package com.org.quoteservice.shared.ws;

import com.rabbitmq.stream.Address;
import com.rabbitmq.stream.AddressResolver;
import com.rabbitmq.stream.ByteCapacity;
import com.rabbitmq.stream.Environment;
import java.time.Duration;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.core.Queue;
import org.springframework.amqp.core.QueueBuilder;
import org.springframework.amqp.rabbit.annotation.EnableRabbit;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.rabbit.stream.config.StreamRabbitListenerContainerFactory;
import org.springframework.rabbit.stream.producer.RabbitStreamTemplate;
import org.springframework.rabbit.stream.support.StreamAdmin;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Configuration
@EnableRabbit
@RequiredArgsConstructor
public class RabbitStreamConfig {

    private final RabbitStreamProperties properties;

    @Bean
    @Primary
    public Environment streamEnvironment() {
        AddressResolver resolver = address -> new Address(properties.host(), address.port());
        log.info("Rabbit stream resolved host={} port={} user={}",
                properties.host(), properties.port(), properties.username());

        return Environment.builder()
                .addressResolver(resolver)
                .host(properties.host())
                .port(properties.port())
                .username(properties.username())
                .password(properties.password())
                .build();
    }

    @Bean
    public StreamAdmin streamAdmin(Environment environment) {
        return new StreamAdmin(environment, spec -> spec.stream(properties.streamName())
                .maxLengthBytes(ByteCapacity.MB(100))
                .maxAge(Duration.ofHours(24))
                .create());
    }

    @Bean
    public Queue quotesStreamQueue() {
        return QueueBuilder.durable(properties.streamName()).stream().build();
    }

    @Bean
    public RabbitStreamTemplate rabbitStreamTemplate(Environment environment) {
        RabbitStreamTemplate template = new RabbitStreamTemplate(environment, properties.streamName());
        template.setMessageConverter(new Jackson2JsonMessageConverter());
        return template;
    }

    @Bean(name = "rabbitListenerContainerFactory")
    public StreamRabbitListenerContainerFactory streamRabbitListenerContainerFactory(Environment environment) {
        StreamRabbitListenerContainerFactory factory = new StreamRabbitListenerContainerFactory(environment);
        factory.setConsumerCustomizer((id, builder) -> builder.name(properties.consumerName()));
        return factory;
    }

}
