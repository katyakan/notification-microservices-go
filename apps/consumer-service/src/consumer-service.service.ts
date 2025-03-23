import { Injectable } from "@nestjs/common";

@Injectable()
export class ConsumerServiceService {
  getHello(): string {
    return "Hello World!";
  }
}
