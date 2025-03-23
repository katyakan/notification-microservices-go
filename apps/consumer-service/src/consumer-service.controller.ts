import { Controller, Get } from "@nestjs/common";
import { ConsumerServiceService } from "./consumer-service.service";

@Controller()
export class ConsumerServiceController {
  constructor(private readonly consumerServiceService: ConsumerServiceService) {}

  @Get()
  getHello(): string {
    return this.consumerServiceService.getHello();
  }
}
