import { Controller, Post, Body, HttpStatus } from "@nestjs/common";
import { ApiTags, ApiOperation, ApiResponse } from "@nestjs/swagger";
import { KafkaService } from "../services/kafka.service";
import { CreateMessageDto } from "../dto/create-message.dto";

@ApiTags("Producer")
@Controller("messages")
export class ProducerController {
  constructor(private readonly kafkaService: KafkaService) {}

  @Post()
  @ApiOperation({ summary: "Send message to Kafka" })
  @ApiResponse({
    status: HttpStatus.CREATED,
    description: "Message has been successfully sent",
  })
  @ApiResponse({
    status: HttpStatus.BAD_REQUEST,
    description: "Invalid message format",
  })
  async sendMessage(@Body() createMessageDto: CreateMessageDto) {
    return this.kafkaService.sendMessage(createMessageDto);
  }
}
