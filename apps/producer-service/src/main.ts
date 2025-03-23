import { NestFactory } from '@nestjs/core';
import { ProducerServiceModule } from './producer-service.module';

async function bootstrap() {
  const app = await NestFactory.create(ProducerServiceModule);
  await app.listen(process.env.port ?? 3000);
}
bootstrap();
