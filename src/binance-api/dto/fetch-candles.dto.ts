import { IsEnum, IsInt, IsOptional, IsString, Max, Min } from 'class-validator';
import { Interval } from '../enums/interval.enum';

export class FetchCandlesDto {
  @IsString()
  symbol: string;

  @IsEnum(Interval)
  interval: Interval;

  @IsOptional()
  @IsInt()
  startTime?: number;

  @IsOptional()
  @IsInt()
  endTime?: number;

  @IsOptional()
  @IsInt()
  @Min(1)
  @Max(1500)
  limit?: number = 500;
}
