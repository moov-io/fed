# FedACH directory file format

**Source:** [achFormat](https://frbservices.org/EPaymentsDirectory/achFormat.html)

| Field Name | Length | Position | Description |
| --- | --- | --- | --- |
| Routing Number | 9 | 1-9 | The institution's routing number |
| Office Code | 1 | 10 | Main office or branch O=main B=branch |
| Servicing FRB Number | 9 | 11-19 | Servicing Fed's main office routing number |
| Record Type Code | 1 | 20 | The code indicating the ABA number to be used to route or send ACH items to the RFI <br/> 0 = Institution is a Federal Reserve Bank <br/> 1 = Send items to customer routing number <br/> 2 = Send items to customer using new routing number field |
| Change Date | 6 | 21-26 | Date of last change to CRF information (MMDDYY) |
| New Routing Number | 9 | 27-35 | Institution's new routing number resulting from a merger or renumber |
| Customer Name | 36 | 36-71 | Commonly used abbreviated name |
| Address | 36 | 72-107 | Delivery address |
| City| 20 | 108-127 | City name in the delivery address |
| State | 2 | 128-129 | State code of the state in the delivery address |
| Zipcode | 5 | 130-134 | Zipcode in the delivery address |
| Zipcode Extension | 4 | 135-138 | Zipcode extension in the delivery address |
| Telephone Area Code | 3 | 139-141 | Area code of the CRF contact telephone number |
| Telephone Prefix Number | 3 | 142-144 | Prefix of the CRF contact telephone number |
| Telephone Suffix Number | 4 | 145-148 | Suffix of the CRF contact telephone number |
| Institution Status Code | 1 | 149 | Code is based on the customers receiver code<br/>1 = Receives Gov/Comm |
| Data View Code | 1 | 150 | 1 = Current view |
| Filler | 5 | 151-155 | Spaces |

# Fedwire directory file format

**Source:** [FedWireFormat](https://frbservices.org/EPaymentsDirectory/fedwireFormat.html)

| Field Name | Length | Columns |
| --- | --- | --- |
| Routing Number | 9 | 1-9 |
| Telegraphic Name | 18 | 10-27 |
| Customer Name | 36 | 28-63 |
| State  | 2 | 64-65 |
| City | 25 | 66-90 |
| Funds transfer status: <br/> Y - Eligible <br/> N - Ineligible | 1 | 91 |
| Funds settlement-only status: <br/> S - Settlement-Only | 1 | 92 |
| Book-Entry Securities transfer status: <br/> Y - Eligible <br/> N - Ineligible | 1 | 93 |
| Date of last revision: YYYYMMDD, or blank | 8 | 94-101 |

