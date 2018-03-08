import { TestBed, inject } from '@angular/core/testing';

import { AddressbookApiClientService } from './addressbook-api-client.service';
import { HttpModule } from '@angular/http';

describe('AddressbookApiClientService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpModule],
      providers: [AddressbookApiClientService]
    });
  });

  it('should be created', inject([AddressbookApiClientService], (service: AddressbookApiClientService) => {
    expect(service).toBeTruthy();
  }));
});
