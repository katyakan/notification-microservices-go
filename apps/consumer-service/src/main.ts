import { NestFactory } from "@nestjs/core";
import { ConsumerServiceModule } from "./consumer-service.module";

async function bootstrap() {
  const app = await NestFactory.create(ConsumerServiceModule);
  await app.listen(process.env.port ?? 3000);
}
bootstrap();
