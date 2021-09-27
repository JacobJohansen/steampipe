package modconfig

// ResourceMetadata ius a struct containing additional data
// about each resource, used to populate the introspection tables
type ResourceMetadata struct {
	ResourceName string `column:"resource_name,text"`
	// mod short name
	ModName          string `column:"mod_name,text"`
	FileName         string `column:"file_name,text"`
	StartLineNumber  int    `column:"start_line_number,integer"`
	EndLineNumber    int    `column:"end_line_number,integer"`
	IsAutoGenerated  bool   `column:"auto_generated,bool"`
	SourceDefinition string `column:"source_definition,text"`
}

// SetMod sets the mod name and mod short name
func (m *ResourceMetadata) SetMod(mod *Mod) {
	// if the mod is the auto-generated default workspace mod, do not save in metadata
	if mod.IsDefaultMod() {
		return
	}
	m.ModName = mod.ShortName
}

// TODO ADD PATH ltree
