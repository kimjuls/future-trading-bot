import { Test, TestingModule } from '@nestjs/testing';
import { BinanceApiService } from './binance-api.service';
import { FetchCandlesDto } from './dto/fetch-candles.dto';
import { Interval } from './enums/interval.enum';
import { lastValueFrom, map } from 'rxjs';
import { HttpModule } from '@nestjs/axios';
import { ConfigModule } from '@nestjs/config';

describe('BinanceApiService', () => {
  let service: BinanceApiService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      imports: [ConfigModule.forRoot(), HttpModule],
      providers: [BinanceApiService],
    }).compile();

    service = module.get<BinanceApiService>(BinanceApiService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('', async () => {
    const dto = new FetchCandlesDto();
    dto.symbol = 'BTCUSDT';
    dto.interval = Interval.OneHour;
    const res = await service.fetchCandles(dto);
    expect(res.status).toEqual(200);
  });
});
