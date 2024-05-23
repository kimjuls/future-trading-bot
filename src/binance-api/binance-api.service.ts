import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { ApiKey } from './classes/api-key';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class BinanceApiService {
  private readonly apiKey: ApiKey;

  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService,
  ) {}

  private configure(): void {
    if (this.apiKey === null) {
      const apiKey = new ApiKey();
      apiKey.access = this.configService.getOrThrow<string>('BINANCE_ACCESS');
      apiKey.secret = this.configService.getOrThrow<string>('BINANCE_SECRET');
    }
  }
}
