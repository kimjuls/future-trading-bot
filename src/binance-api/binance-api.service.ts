import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { ApiKey } from './classes/api-key';
import { ConfigService } from '@nestjs/config';
import { createHmac } from 'crypto';
import { FetchCandlesDto } from './dto/fetch-candles.dto';
import { RawAxiosRequestConfig } from 'axios';
import { lastValueFrom, map } from 'rxjs';

@Injectable()
export class BinanceApiService {
  private apiKey: ApiKey = null;
  private readonly baseUrl = 'https://fapi.binance.com';

  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService,
  ) {
    this.configure();
  }

  private configure(): void {
    if (this.apiKey === null) {
      const apiKey = new ApiKey();
      apiKey.access = this.configService.getOrThrow<string>('BINANCE_ACCESS');
      apiKey.secret = this.configService.getOrThrow<string>('BINANCE_SECRET');
      this.apiKey = apiKey;
    }
  }

  private sign(params?: any) {
    this.configure();
    const queryString = new URLSearchParams(params).toString();
    const signature = createHmac('sha256', this.apiKey.secret)
      .update(queryString)
      .digest('hex');
    return signature;
  }

  async fetchCandles(dto: FetchCandlesDto) {
    const url = '/fapi/v1/klines';
    const config: RawAxiosRequestConfig = {
      method: 'GET',
      baseURL: this.baseUrl,
      url,
      params: dto,
    };
    const response = await lastValueFrom(this.httpService.request(config));
    return response;
  }
}
