import { IsString, IsNotEmpty, IsObject } from "class-validator";
import { ApiProperty } from "@nestjs/swagger";

export class CreateMessageDto {
  @ApiProperty({
    description: "Type of the message",
    example: "notification",
  })
  @IsString()
  @IsNotEmpty()
  type: string;

  @ApiProperty({
    description: "Message payload",
    example: { text: "Hello World", chatId: 123456 },
  })
  @IsObject()
  @IsNotEmpty()
  payload: Record<string, any>;
}
