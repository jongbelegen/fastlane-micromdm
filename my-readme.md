Mogelijk dat we dus eerst hexnode mdm vanuit apple business manager moeten verwijderen en dan pas micromdm moeten toevoegen.

Reassignment Oddities
There appears to be some oddities when assigning device serials from one MDM server to another in the ABM/ASM/DEP portal. When re-assigning a serial number it appears that instead of generating an "added" event (which MicroMDM will use to auto-assign serial numbers) Apple instead generates a "modified" event — even if that MDM server has never seen the device before. Because these are "modified" events the auto-assigner won't work on those serials. This appears to be some issue with the ABM/ASM/DEP portal.

As a workaround it's been discovered that if you add a step in between your assignment on the Apple portal you can get it to successfully generate your "added" events for auto-assignment. This involves either unassigning the serial numbers from DEP or re-assigning your serial numbers to another

// Commands zijn altijd async en worden pas uitgevoerd wanneer de ipad incheckt bij de mdm server
We kunnen inchecken versnellen soms door een notificatie te sturen.
Stel de ipad stond uit, en de mdm stuurt een notifiactie dan heeft de ipad mogelijk nooit ingecheckt.
Dan moeten we hem opnieuw "nudgen" zodat de queue naar beneden gaat

