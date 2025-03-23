import { Controller, Get } from "@nestjs/common";
import { ApiTags, ApiOperation } from "@nestjs/swagger";

@ApiTags("Consumer")
@Controller("consumer")
export class ConsumerController {
  @Get("health")
  @ApiOperation({ summary: "Get consumer service health status" })
  getHealth() {
    return { status: "ok" };
  }
}
