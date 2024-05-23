export class ApiKey {
  private _access: string;
  private _secret: string;
  constructor() {}
  public get access(): string {
    return this._access;
  }
  public set access(value: string) {
    this._access = value;
  }
  public get secret(): string {
    return this._secret;
  }
  public set secret(value: string) {
    this._secret = value;
  }
}
