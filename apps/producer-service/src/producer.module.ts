import { Module } from "@nestjs/common";
import { ProducerController } from "./controllers/producer.controller";
import { KafkaService } from "./services/kafka.service";
import { ConfigModule } from '@nestjs/config';
@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
  ],
  controllers: [ProducerController],
  providers: [KafkaService],
})
export class ProducerModule {}
