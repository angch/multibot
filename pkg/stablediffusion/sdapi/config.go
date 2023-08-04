package sdapi

import (
	"bytes"
	"encoding/json"
)

type intNum int

func (c *intNum) UnmarshalJSON(data []byte) error {
	i := 0
	f := 0.0
	if err := json.Unmarshal(data, &i); err == nil {
		*c = intNum(i)
		return nil
	}
	if err := json.Unmarshal(data, &f); err == nil {
		*c = intNum(f)
		return nil
	}
	return nil
}

type Config struct {
	SamplesSave                        bool        `json:"samples_save"`
	SamplesFormat                      string      `json:"samples_format"`
	SamplesFilenamePattern             string      `json:"samples_filename_pattern"`
	SaveImagesAddNumber                bool        `json:"save_images_add_number"`
	GridSave                           bool        `json:"grid_save"`
	GridFormat                         string      `json:"grid_format"`
	GridExtendedFilename               bool        `json:"grid_extended_filename"`
	GridOnlyIfMultiple                 bool        `json:"grid_only_if_multiple"`
	GridPreventEmptySpots              bool        `json:"grid_prevent_empty_spots"`
	NRows                              intNum      `json:"n_rows"`
	EnablePnginfo                      bool        `json:"enable_pnginfo"`
	SaveTxt                            bool        `json:"save_txt"`
	SaveImagesBeforeFaceRestoration    bool        `json:"save_images_before_face_restoration"`
	SaveImagesBeforeHighresFix         bool        `json:"save_images_before_highres_fix"`
	SaveImagesBeforeColorCorrection    bool        `json:"save_images_before_color_correction"`
	JpegQuality                        intNum      `json:"jpeg_quality"`
	ExportFor4Chan                     bool        `json:"export_for_4chan"`
	UseOriginalNameBatch               bool        `json:"use_original_name_batch"`
	UseUpscalerNameAsSuffix            bool        `json:"use_upscaler_name_as_suffix"`
	SaveSelectedOnly                   bool        `json:"save_selected_only"`
	DoNotAddWatermark                  bool        `json:"do_not_add_watermark"`
	TempDir                            string      `json:"temp_dir"`
	CleanTempDirAtStart                bool        `json:"clean_temp_dir_at_start"`
	OutdirSamples                      string      `json:"outdir_samples"`
	OutdirTxt2ImgSamples               string      `json:"outdir_txt2img_samples"`
	OutdirImg2ImgSamples               string      `json:"outdir_img2img_samples"`
	OutdirExtrasSamples                string      `json:"outdir_extras_samples"`
	OutdirGrids                        string      `json:"outdir_grids"`
	OutdirTxt2ImgGrids                 string      `json:"outdir_txt2img_grids"`
	OutdirImg2ImgGrids                 string      `json:"outdir_img2img_grids"`
	OutdirSave                         string      `json:"outdir_save"`
	SaveToDirs                         bool        `json:"save_to_dirs"`
	GridSaveToDirs                     bool        `json:"grid_save_to_dirs"`
	UseSaveToDirsForUI                 bool        `json:"use_save_to_dirs_for_ui"`
	DirectoriesFilenamePattern         string      `json:"directories_filename_pattern"`
	DirectoriesMaxPromptWords          intNum      `json:"directories_max_prompt_words"`
	ESRGANTile                         intNum      `json:"ESRGAN_tile"`
	ESRGANTileOverlap                  intNum      `json:"ESRGAN_tile_overlap"`
	RealesrganEnabledModels            []string    `json:"realesrgan_enabled_models"`
	UpscalerForImg2Img                 *string     `json:"upscaler_for_img2img"`
	UseScaleLatentForHiresFix          bool        `json:"use_scale_latent_for_hires_fix"`
	LdsrSteps                          intNum      `json:"ldsr_steps"`
	LdsrCached                         bool        `json:"ldsr_cached"`
	SWINTile                           intNum      `json:"SWIN_tile"`
	SWINTileOverlap                    intNum      `json:"SWIN_tile_overlap"`
	FaceRestorationModel               string      `json:"face_restoration_model"`
	CodeFormerWeight                   float32     `json:"code_former_weight"`
	FaceRestorationUnload              bool        `json:"face_restoration_unload"`
	MemmonPollRate                     float32     `json:"memmon_poll_rate"`
	SamplesLogStdout                   bool        `json:"samples_log_stdout"`
	MultipleTqdm                       bool        `json:"multiple_tqdm"`
	UnloadModelsWhenTraining           bool        `json:"unload_models_when_training"`
	PinMemory                          bool        `json:"pin_memory"`
	SaveOptimizerState                 bool        `json:"save_optimizer_state"`
	DatasetFilenameWordRegex           string      `json:"dataset_filename_word_regex"`
	DatasetFilenameJoinString          string      `json:"dataset_filename_join_string"`
	TrainingImageRepeatsPerEpoch       intNum      `json:"training_image_repeats_per_epoch"`
	TrainingWriteCsvEvery              float32     `json:"training_write_csv_every"`
	TrainingXattentionOptimizations    bool        `json:"training_xattention_optimizations"`
	SdModelCheckpoint                  string      `json:"sd_model_checkpoint"`
	SdCheckpointCache                  intNum      `json:"sd_checkpoint_cache"`
	SdVaeCheckpointCache               intNum      `json:"sd_vae_checkpoint_cache"`
	SdVae                              string      `json:"sd_vae"`
	SdVaeAsDefault                     bool        `json:"sd_vae_as_default"`
	SdHypernetwork                     string      `json:"sd_hypernetwork"`
	SdHypernetworkStrength             float32     `json:"sd_hypernetwork_strength"`
	InpaintingMaskWeight               float32     `json:"inpainting_mask_weight"`
	InitialNoiseMultiplier             json.Number `json:"initial_noise_multiplier"`
	Img2ImgColorCorrection             bool        `json:"img2img_color_correction"`
	Img2ImgFixSteps                    bool        `json:"img2img_fix_steps"`
	Img2ImgBackgroundColor             string      `json:"img2img_background_color"`
	EnableQuantization                 bool        `json:"enable_quantization"`
	EnableEmphasis                     bool        `json:"enable_emphasis"`
	UseOldEmphasisImplementation       bool        `json:"use_old_emphasis_implementation"`
	EnableBatchSeeds                   bool        `json:"enable_batch_seeds"`
	CommaPaddingBacktrack              intNum      `json:"comma_padding_backtrack"`
	CLIPStopAtLastLayers               intNum      `json:"CLIP_stop_at_last_layers"`
	RandomArtistCategories             []string    `json:"random_artist_categories"`
	InterrogateKeepModelsInMemory      bool        `json:"interrogate_keep_models_in_memory"`
	InterrogateUseBuiltinArtists       bool        `json:"interrogate_use_builtin_artists"`
	InterrogateReturnRanks             bool        `json:"interrogate_return_ranks"`
	InterrogateClipNumBeams            intNum      `json:"interrogate_clip_num_beams"`
	InterrogateClipMinLength           intNum      `json:"interrogate_clip_min_length"`
	InterrogateClipMaxLength           intNum      `json:"interrogate_clip_max_length"`
	InterrogateClipDictLimit           intNum      `json:"interrogate_clip_dict_limit"`
	InterrogateDeepbooruScoreThreshold float32     `json:"interrogate_deepbooru_score_threshold"`
	DeepbooruSortAlpha                 bool        `json:"deepbooru_sort_alpha"`
	DeepbooruUseSpaces                 bool        `json:"deepbooru_use_spaces"`
	DeepbooruEscape                    bool        `json:"deepbooru_escape"`
	DeepbooruFilterTags                string      `json:"deepbooru_filter_tags"`
	ShowProgressbar                    bool        `json:"show_progressbar"`
	ShowProgressEveryNSteps            intNum      `json:"show_progress_every_n_steps"`
	ShowProgressType                   string      `json:"show_progress_type"`
	ShowProgressGrid                   bool        `json:"show_progress_grid"`
	ReturnGrid                         bool        `json:"return_grid"`
	DoNotShowImages                    bool        `json:"do_not_show_images"`
	AddModelHashToInfo                 bool        `json:"add_model_hash_to_info"`
	AddModelNameToInfo                 bool        `json:"add_model_name_to_info"`
	DisableWeightsAutoSwap             bool        `json:"disable_weights_auto_swap"`
	SendSeed                           bool        `json:"send_seed"`
	SendSize                           bool        `json:"send_size"`
	Font                               string      `json:"font"`
	JsModalLightbox                    bool        `json:"js_modal_lightbox"`
	JsModalLightboxInitiallyZoomed     bool        `json:"js_modal_lightbox_initially_zoomed"`
	ShowProgressInTitle                bool        `json:"show_progress_in_title"`
	SamplersInDropdown                 bool        `json:"samplers_in_dropdown"`
	Quicksettings                      string      `json:"quicksettings"`
	Localization                       string      `json:"localization"`
	HideSamplers                       []string    `json:"hide_samplers"`
	EtaDdim                            intNum      `json:"eta_ddim"`
	EtaAncestral                       intNum      `json:"eta_ancestral"`
	DdimDiscretize                     string      `json:"ddim_discretize"`
	SChurn                             intNum      `json:"s_churn"`
	STmin                              intNum      `json:"s_tmin"`
	SNoise                             intNum      `json:"s_noise"`
	EtaNoiseSeedDelta                  intNum      `json:"eta_noise_seed_delta"`
	DisabledExtensions                 []string    `json:"disabled_extensions"`
}

func (c *Config) IoReader() *bytes.Buffer {
	j, err := json.Marshal(c)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(j)
}
